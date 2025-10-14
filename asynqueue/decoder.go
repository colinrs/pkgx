package asynqueue

import "encoding/json"

// Decoder 解码器抽象（可换 msgpack/proto）
type Decoder func([]byte, any) error

func JSONDecoder(b []byte, v any) error { return json.Unmarshal(b, v) }
