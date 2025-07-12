package main

import "os"

func service2key(service string) (string, string) {
	// 这里可以根据 service 的值返回不同的 API 密钥和端点
	switch service {
	case "chutes":
		return os.Getenv("CHUTES_API_TOKEN"), "https://llm.chutes.ai/v1/chat/completions"
	case "chutes-hidream": // text2pic
		return os.Getenv("CHUTES_API_TOKEN"), "https://chutes-hidream.chutes.ai/generate"
	case "chutes-chroma": // text2pic
		return os.Getenv("CHUTES_API_TOKEN"), "https://chutes-chroma.chutes.ai/generate"
	case "chutes-stable-flow": // text2pic
		return os.Getenv("CHUTES_API_TOKEN"), "https://chutes-stable-flow.chutes.ai/generate"
	case "chutes-infiniteyou": // text2pic
		return os.Getenv("CHUTES_API_TOKEN"), "https://chutes-infiniteyou.chutes.ai/generate"
	case "groq":
		return os.Getenv("GROQ_API_KEY"), "https://api.groq.com/openai/v1/chat/completions"
	case "huawei-ds-v3":
		return os.Getenv("HUAWEI_API_KEY"), "https://maas-cn-southwest-2.modelarts-maas.com/v1/infers/271c9332-4aa6-4ff5-95b3-0cf8bd94c394/v1/chat/completions"
	case "huawei-ds-r1":
		return os.Getenv("HUAWEI_API_KEY"), "https://maas-cn-southwest-2.modelarts-maas.com/v1/infers/8a062fd4-7367-4ab4-a936-5eeb8fb821c4/v1/chat/completions"
	default:
		return os.Getenv("GROQ_API_KEY"), "https://api.groq.com/openai/v1/chat/completions"
	}
}
