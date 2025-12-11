import logging
import os
from dataclasses import asdict, dataclass, field
from typing import Any, Literal

import requests
from tenacity import (
    before_sleep_log,
    retry,
    retry_if_not_exception_type,
    stop_after_attempt,
    wait_exponential,
)

from minisweagent.models import GLOBAL_MODEL_STATS

logger = logging.getLogger("stk_model")


@dataclass
class StkModelConfig:
    model_name: str  # Agent ID (required)
    model_kwargs: dict[str, Any] = field(default_factory=dict)
    client_id: str | None = os.getenv("STK_CLIENT_ID")
    client_secret: str | None = os.getenv("STK_CLIENT_SECRET")
    realm: str | None = os.getenv("STK_REALM")
    cost_tracking: Literal["default", "ignore_errors"] = os.getenv("MSWEA_COST_TRACKING", "default")


class StkModel:
    def __init__(self, *, config_class: type = StkModelConfig, **kwargs):
        self.config = config_class(**kwargs)
        self.cost = 0.0
        self.n_calls = 0
        self._token = None

        if not self.config.model_name:
            raise ValueError("Agent ID (model_name) is required.")
        if not self.config.client_id:
            raise ValueError("STK_CLIENT_ID is required.")
        if not self.config.client_secret:
            raise ValueError("STK_CLIENT_SECRET is required.")
        if not self.config.realm:
            raise ValueError("STK_REALM is required.")

        self.conversation_id = None

    def _get_token(self) -> str:
        # Simple token retrieval, no caching expiration check for now (can be added if needed)
        # In a real app, we should check expiration. For now, we'll just fetch if we don't have one
        # or if a request fails (handled in query retry logic potentially, but let's just fetch fresh for now to be safe or cache simply)
        # The user's script fetches it every time. Let's fetch it if not present.
        # Ideally we should refresh on 401.
        
        url = f"https://idm.stackspot.com/{self.config.realm}/oidc/oauth/token"
        payload = {
            "grant_type": "client_credentials",
            "client_id": self.config.client_id,
            "client_secret": self.config.client_secret,
        }
        headers = {"Content-Type": "application/x-www-form-urlencoded"}
        
        response = requests.post(url, data=payload, headers=headers)
        response.raise_for_status()
        return response.json()["access_token"]

    @retry(
        stop=stop_after_attempt(int(os.getenv("MSWEA_MODEL_RETRY_STOP_AFTER_ATTEMPT", "10"))),
        wait=wait_exponential(multiplier=1, min=4, max=60),
        before_sleep=before_sleep_log(logger, logging.WARNING),
        retry=retry_if_not_exception_type((KeyboardInterrupt, ValueError)),
    )
    def _query(self, messages: list[dict[str, str]], **kwargs):
        # Extract the last user message as the prompt
        user_prompt = ""
        for msg in reversed(messages):
            if msg["role"] == "user":
                user_prompt = msg["content"]
                break
        
        if not user_prompt:
             # Fallback or empty?
             pass

        if not self._token:
            self._token = self._get_token()

        url = f"https://genai-inference-app.stackspot.com/v1/agent/{self.config.model_name}/chat"
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self._token}",
        }
        
        data = {
            "streaming": False,
            "user_prompt": user_prompt,
            "stackspot_knowledge": False,
            "return_ks_in_response": False,
            "use_conversation": True,
        }
        
        if self.conversation_id:
            data["conversation_id"] = self.conversation_id

        # Merge extra kwargs
        data.update(self.config.model_kwargs)
        data.update(kwargs)

        try:
            response = requests.post(url, json=data, headers=headers)
            if response.status_code == 401:
                # Token might be expired, refresh and retry
                logger.warning("Token expired, refreshing...")
                self._token = self._get_token()
                headers["Authorization"] = f"Bearer {self._token}"
                response = requests.post(url, json=data, headers=headers)
            
            response.raise_for_status()
            response_json = response.json()
            
            # Update conversation_id if present
            if "conversation_id" in response_json:
                self.conversation_id = response_json["conversation_id"]
                
            return response_json
        except requests.RequestException as e:
            logger.error(f"Request failed: {e}")
            raise

    def query(self, messages: list[dict[str, str]], **kwargs) -> dict:
        response_json = self._query(messages, **kwargs)
        
        # Parse response based on user provided format:
        # { "message": "...", "stop_reason": "...", "conversation_id": "..." }
        content = response_json.get("message") or ""
        
        if not content:
            # Fallback if message is null or empty, though it shouldn't be for a valid response
            content = str(response_json)

        cost = 0.0
        self.n_calls += 1
        self.cost += cost
        GLOBAL_MODEL_STATS.add(cost)
        
        return {
            "content": content,
            "extra": {
                "response": response_json,
                "cost": cost,
            },
        }

    def get_template_vars(self) -> dict[str, Any]:
        return asdict(self.config) | {"n_model_calls": self.n_calls, "model_cost": self.cost}
