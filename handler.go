package main

import "sync"

type SetCmd struct {
	HsetMap map[string]string
	MU      sync.RWMutex
}

var SETs = SetCmd{
	HsetMap: make(map[string]string),
	MU:      sync.RWMutex{},
}

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR Wrong number of arguments"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETs.MU.Lock()
	defer SETs.MU.Unlock()
	SETs.HsetMap[key] = value

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR Wrong number of arguments"}
	}

	key := args[0].bulk

	SETs.MU.RLock()
	defer SETs.MU.RUnlock()

	value, ok := SETs.HsetMap[key]
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

type HSetCmd struct {
	HsetMap map[string]map[string]string
	MU      sync.RWMutex
}

var HSETs HSetCmd = HSetCmd{
	HsetMap: make(map[string]map[string]string),
	MU:      sync.RWMutex{},
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR Wrong number of arguments"}
	}

	category := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETs.MU.Lock()
	defer HSETs.MU.Unlock()

	if _, ok := HSETs.HsetMap[category]; !ok {
		HSETs.HsetMap[category] = make(map[string]string)
	}

	HSETs.HsetMap[category][key] = value

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR Wrong number of arguments"}
	}

	category := args[0].bulk
	key := args[1].bulk

	HSETs.MU.RLock()
	defer HSETs.MU.RUnlock()

	value, ok := HSETs.HsetMap[category][key]
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR Wrong number of arguments"}
	}

	category := args[0].bulk

	HSETs.MU.RLock()
	defer HSETs.MU.RUnlock()

	value, ok := HSETs.HsetMap[category]
	if !ok {
		return Value{typ: "null"}
	}

	out := Value{typ: "array"}
	for key, val := range value {
		out.array = append(out.array, Value{typ: "bulk", bulk: key})
		out.array = append(out.array, Value{typ: "bulk", bulk: val})
	}

	return out
}
