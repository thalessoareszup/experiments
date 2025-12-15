import { useEffect, useState, useCallback, useRef } from 'react';
import type { Plan, Step } from '../types';

interface UseWebSocketOptions {
  onPlanCreated?: (plan: Plan) => void;
  onPlanUpdated?: (plan: Plan) => void;
  onPlanDeleted?: (id: string) => void;
  onStepCreated?: (step: Step) => void;
  onStepUpdated?: (step: Step) => void;
  onStepDeleted?: (id: string) => void;
  reconnectInterval?: number;
  maxReconnectInterval?: number;
}

interface WebSocketMessage {
  type: string;
  data: Plan | Step | { id: string };
  timestamp: string;
  id: string;
}

type ConnectionState = 'connecting' | 'connected' | 'disconnected' | 'error';

// LRU cache for message deduplication
class MessageIDCache {
  private cache: Set<string>;
  private readonly maxSize = 1000;

  has(id: string): boolean {
    return this.cache.has(id);
  }

  add(id: string) {
    if (this.cache.size >= this.maxSize) {
      // Remove oldest entry (this is a simple implementation)
      const first = this.cache.values().next().value;
      if (first) this.cache.delete(first);
    }
    this.cache.add(id);
  }

  constructor() {
    this.cache = new Set();
  }
}

export function useWebSocket(url: string, options: UseWebSocketOptions = {}) {
  const [connectionState, setConnectionState] = useState<ConnectionState>('disconnected');
  const [error, setError] = useState<string | null>(null);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const heartbeatTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const messageIDCacheRef = useRef(new MessageIDCache());

  const reconnectInterval = options.reconnectInterval ?? 1000;
  const maxReconnectInterval = options.maxReconnectInterval ?? 30000;

  const handleMessage = useCallback(
    (message: WebSocketMessage) => {
      // Deduplicate messages
      if (messageIDCacheRef.current.has(message.id)) {
        console.debug('Duplicate message ignored:', message.id);
        return;
      }
      messageIDCacheRef.current.add(message.id);

      // Handle different event types
      switch (message.type) {
        case 'plan:created':
          options.onPlanCreated?.(message.data as Plan);
          break;
        case 'plan:updated':
          options.onPlanUpdated?.(message.data as Plan);
          break;
        case 'plan:deleted':
          options.onPlanDeleted?.((message.data as { id: string }).id);
          break;
        case 'step:created':
          options.onStepCreated?.(message.data as Step);
          break;
        case 'step:updated':
          options.onStepUpdated?.(message.data as Step);
          break;
        case 'step:deleted':
          options.onStepDeleted?.((message.data as { id: string }).id);
          break;
        case 'connected':
          setConnectionState('connected');
          setError(null);
          break;
        default:
          console.debug('Unknown message type:', message.type);
      }
    },
    [options],
  );

  const resetHeartbeat = useCallback(() => {
    if (heartbeatTimeoutRef.current) {
      clearTimeout(heartbeatTimeoutRef.current);
    }
    // If no message received in 45 seconds, reconnect
    heartbeatTimeoutRef.current = setTimeout(() => {
      console.warn('Heartbeat timeout, reconnecting...');
      if (wsRef.current) {
        wsRef.current.close();
      }
    }, 45000);
  }, []);

  const connect = useCallback(() => {
    // Convert http/https to ws/wss
    const wsUrl = url
      .replace(/^http:/, 'ws:')
      .replace(/^https:/, 'wss:');

    console.log('WebSocket connecting to:', wsUrl);
    setConnectionState('connecting');

    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log('WebSocket connected');
      setConnectionState('connected');
      setError(null);
      reconnectAttemptsRef.current = 0;
      resetHeartbeat();
    };

    ws.onmessage = (event) => {
      resetHeartbeat();
      try {
        const message = JSON.parse(event.data) as WebSocketMessage;
        handleMessage(message);
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err);
        setError('Failed to parse message from server');
      }
    };

    ws.onerror = (event) => {
      console.error('WebSocket error:', event);
      setConnectionState('error');
      setError('WebSocket connection error');
    };

    ws.onclose = () => {
      console.log('WebSocket closed');
      setConnectionState('disconnected');

      // Auto-reconnect with exponential backoff
      if (reconnectAttemptsRef.current < 10) {
        const delay = Math.min(
          reconnectInterval * Math.pow(2, reconnectAttemptsRef.current),
          maxReconnectInterval,
        );
        console.log(`Reconnecting in ${delay}ms (attempt ${reconnectAttemptsRef.current + 1})`);
        reconnectAttemptsRef.current++;

        reconnectTimeoutRef.current = setTimeout(() => {
          connect();
        }, delay);
      } else {
        setError('Max reconnection attempts reached');
      }
    };

    wsRef.current = ws;

    return ws;
  }, [url, handleMessage, resetHeartbeat, reconnectInterval, maxReconnectInterval]);

  const reconnect = useCallback(() => {
    console.log('Manual reconnect requested');
    if (wsRef.current) {
      wsRef.current.close();
    }
    reconnectAttemptsRef.current = 0;
    connect();
  }, [connect]);

  useEffect(() => {
    connect();

    return () => {
      // Cleanup
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (heartbeatTimeoutRef.current) {
        clearTimeout(heartbeatTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect]);

  return {
    connectionState,
    error,
    reconnect,
  };
}
