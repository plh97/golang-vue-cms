// src/event.ts
// Simple Event Emitter for global app events (unlogin, error, etc.)
export type EventHandler = (payload?: any) => void
const eventBus: Record<string, EventHandler[]> = {}

export function on(event: string, handler: EventHandler) {
  if (!eventBus[event]) eventBus[event] = []
  eventBus[event].push(handler)
}

export function off(event: string, handler: EventHandler) {
  if (!eventBus[event]) return
  eventBus[event] = eventBus[event].filter(h => h !== handler)
}

export function emit(event: string, payload?: any) {
  (eventBus[event] || []).forEach(h => h(payload))
}
