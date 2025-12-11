"""Utilities for displaying model metrics (cost or tokens)."""


def get_metrics_display(model) -> str:
    """Return formatted metrics: tokens or cost based on model type.

    Args:
        model: The model instance

    Returns:
        Formatted string like "$1.23" or "168 tokens (123 in / 45 out)"
    """
    # Check if model has token tracking
    if (
        hasattr(model, "tokens_input")
        and hasattr(model, "tokens_output")
        and model.tokens_input is not None
        and model.tokens_output is not None
    ):
        total = model.tokens_input + model.tokens_output
        return f"{total} tokens ({model.tokens_input} in / {model.tokens_output} out)"
    else:
        return f"${model.cost:.2f}"
