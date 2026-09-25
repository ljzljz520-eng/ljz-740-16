# Stable Diffusion Go Binding (PureGo)

使用 [purego](https://github.com/ebitengine/purego) 实现的 `stable-diffusion.cpp` 纯 Go 语言绑定。  
**无需 CGO**，支持 macOS、Linux 和 Windows 三大平台。

## ✨ 核心亮点

| 特性 | 说明 |
|------|------|
| **零 CGO** | 纯 Go 实现，通过 `purego` 直接调用 C 动态库 |
| **三平台** | macOS `.dylib` / Linux `.so` / Windows `.dll` 均支持 |
| **功能完整** | 1:1 覆盖 `stable-diffusion.h` 中全部 38 个导出函数 |
| **生产就绪** | 示例服务含并发控制、超时处理、PNG 编码、优雅退出 |

---

## 📦 模块路径

```
github.com/example/stablediffusion           # 高层 API
github.com/example/stablediffusion/bindings  # 底层 purego 绑定
```

---

## 🚀 快速上手

### 1. 运行单元测试（无需 C++ 环境）

测试会自动使用 Mock 实现，无需动态库：

```bash
go test -v ./...
```

### 2. macOS / Linux 原生构建

**A. 编译动态库**
```bash
cd stable-diffusion.cpp
mkdir -p build && cd build
cmake .. -DSD_BUILD_SHARED_LIBS=ON -DCMAKE_BUILD_TYPE=Release
cmake --build . --config Release -j
```

**B. 启动服务**
```bash
cd ../..
# macOS
export SD_LIB_PATH=$(pwd)/stable-diffusion.cpp/build/bin/libstable-diffusion.dylib
# Linux
# export SD_LIB_PATH=$(pwd)/stable-diffusion.cpp/build/bin/libstable-diffusion.so

go run examples/server/server.go
```

### 3. Windows 原生构建

**A. 前置条件**
- [CMake](https://cmake.org/download/)
- [Visual Studio Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/) (含 MSVC)

**B. 编译 DLL**
```powershell
cd stable-diffusion.cpp
mkdir build; cd build
cmake .. -DSD_BUILD_SHARED_LIBS=ON -DCMAKE_BUILD_TYPE=Release
cmake --build . --config Release
```
编译产物路径：`stable-diffusion.cpp/build/bin/Release/stable-diffusion.dll`

> **purego 在 Windows 上的工作原理**：`purego` 底层在 Windows 上使用 `syscall.LoadLibrary` / `GetProcAddress` 加载 DLL，与 Unix 上的 `dlopen`/`dlsym` 等价，完全兼容。

**C. 启动服务**
```powershell
cd ../..
$env:SD_LIB_PATH="$(Get-Location)\stable-diffusion.cpp\build\bin\Release\stable-diffusion.dll"
go run examples/server/server.go
```

### 4. Docker 部署（Linux 环境）

```bash
docker compose up --build -d
curl http://localhost:8080/health
```

> **macOS Docker 说明**：Dockerfile 使用 `debian:bookworm-slim` 基础镜像，面向 Linux 服务器部署。在 macOS 上 Docker Desktop 通过虚拟机运行，无法透传 Metal GPU，性能大幅下降。macOS 推荐使用上述 **原生构建** 方式。

---

## 📋 API 概览

### 基础用法

```go
package main

import (
    "log"
    "github.com/example/stablediffusion"
    "github.com/example/stablediffusion/bindings"
)

func main() {
    // 创建上下文
    options := stablediffusion.DefaultContextOptions("models/sd1.5.gguf")
    options.NThreads = 8

    ctx, err := stablediffusion.NewContext(options)
    if err != nil {
        log.Fatalf("Failed to create context: %v", err)
    }
    defer ctx.Free()

    // 图像生成
    cfg := stablediffusion.GenerationConfig{
        Prompt:     "A cyberpunk city, 8k wallpaper",
        Width:      512,
        Height:     512,
        BatchCount: 1,
        Sampler: stablediffusion.SamplerConfig{
            Steps:     20,
            TxtCfg:    7.0,
            Method:    bindings.EULER_A_SAMPLE_METHOD,
            Scheduler: bindings.KARRAS_SCHEDULER,
        },
    }

    images, err := ctx.GenerateImage(cfg)
    if err != nil {
        log.Fatalf("Failed: %v", err)
    }
    log.Printf("Generated %d images", len(images))
}
```

### 完整 API 清单

对比 `stable-diffusion.h` 全部 38 个导出函数，以下是绑定覆盖状态：

| C 函数 | Go 绑定 (bindings) | 高层 API (stablediffusion) |
|--------|-------|-----------|
| `sd_set_log_callback` | `SetLogCallback` | `SetLogCallback` |
| `sd_set_progress_callback` | `SetProgressCallback` | `SetProgressCallback` |
| `sd_set_preview_callback` | `SetPreviewCallback` | `SetPreviewCallback` |
| `sd_get_num_physical_cores` | `GetNumPhysicalCores` | `GetNumPhysicalCores` |
| `sd_get_system_info` | `GetSystemInfo` | `GetSystemInfo` |
| `sd_type_name` | `GetTypeName` | `GetTypeName` |
| `str_to_sd_type` | `StrToSdType` | `StrToSdType` |
| `sd_rng_type_name` | `GetRngTypeName` | `GetRngTypeName` |
| `str_to_rng_type` | `StrToRngType` | `StrToRngType` |
| `sd_sample_method_name` | `GetSampleMethodName` | `GetSampleMethodName` |
| `str_to_sample_method` | `StrToSampleMethod` | `StrToSampleMethod` |
| `sd_scheduler_name` | `GetSchedulerName` | `GetSchedulerName` |
| `str_to_scheduler` | `StrToScheduler` | `StrToScheduler` |
| `sd_prediction_name` | `GetPredictionName` | `GetPredictionName` |
| `str_to_prediction` | `StrToPrediction` | `StrToPrediction` |
| `sd_preview_name` | `GetPreviewName` | `GetPreviewName` |
| `str_to_preview` | `StrToPreview` | `StrToPreview` |
| `sd_lora_apply_mode_name` | `GetLoraApplyModeName` | `GetLoraApplyModeName` |
| `str_to_lora_apply_mode` | `StrToLoraApplyMode` | `StrToLoraApplyMode` |
| `sd_cache_params_init` | `SdCacheParamsInit` | ✅ 内部使用 |
| `sd_ctx_params_init` | `SdCtxParamsInit` | ✅ 内部使用 |
| `sd_ctx_params_to_str` | `SdCtxParamsToStr` | ✅ 调试用 |
| `new_sd_ctx` | `CreateSdCtx` | `NewContext` |
| `free_sd_ctx` | `FreeSdCtx` | `Context.Free` |
| `sd_sample_params_init` | `SdSampleParamsInit` | ✅ 内部使用 |
| `sd_sample_params_to_str` | `SdSampleParamsToStr` | ✅ 调试用 |
| `sd_get_default_sample_method` | `GetDefaultSampleMethod` | `Context.GetDefaultSampleMethod` |
| `sd_get_default_scheduler` | `GetDefaultScheduler` | `Context.GetDefaultScheduler` |
| `sd_img_gen_params_init` | `SdImgGenParamsInit` | ✅ 内部使用 |
| `sd_img_gen_params_to_str` | `SdImgGenParamsToStr` | ✅ 调试用 |
| `generate_image` | `GenerateImage` | `Context.GenerateImage` |
| `sd_vid_gen_params_init` | `SdVidGenParamsInit` | ✅ 内部使用 |
| `generate_video` | `GenerateVideo` | `Context.GenerateVideo` |
| `new_upscaler_ctx` | `CreateUpscalerCtx` | `NewUpscaler` |
| `free_upscaler_ctx` | `FreeUpscalerCtx` | `Upscaler.Free` |
| `upscale` | `Upscale` | `Upscaler.Upscale` |
| `get_upscale_factor` | `GetUpscaleFactor` | `Upscaler.GetUpscaleFactor` |
| `convert` | `Convert` | `ConvertModel` |
| `preprocess_canny` | `PreprocessCanny` | `PreprocessCanny` |
| `sd_commit` | `GetCommit` | `GetCommit` |
| `sd_version` | `GetVersion` | `GetVersion` |

---

## 🎛 支持的采样器与调度器

采样器（`SampleMethod`）与调度器（`Scheduler`）已整理为独立的 Go 枚举
（见 `bindings/sampler.go`），常量值与 `stable-diffusion.h` 中的
`enum sample_method_t` / `enum scheduler_t` 严格一一对应。

### 采样器（SampleMethod）

| 字符串名称 | Go 常量 | 值 |
|-----------|---------|----|
| `euler` | `bindings.EULER_SAMPLE_METHOD` | 0 |
| `euler_a` | `bindings.EULER_A_SAMPLE_METHOD` | 1 |
| `heun` | `bindings.HEUN_SAMPLE_METHOD` | 2 |
| `dpm2` | `bindings.DPM2_SAMPLE_METHOD` | 3 |
| `dpm++2s_a` | `bindings.DPMPP2S_A_SAMPLE_METHOD` | 4 |
| `dpm++2m` | `bindings.DPMPP2M_SAMPLE_METHOD` | 5 |
| `dpm++2mv2` | `bindings.DPMPP2Mv2_SAMPLE_METHOD` | 6 |
| `ipndm` | `bindings.IPNDM_SAMPLE_METHOD` | 7 |
| `ipndm_v` | `bindings.IPNDM_V_SAMPLE_METHOD` | 8 |
| `lcm` | `bindings.LCM_SAMPLE_METHOD` | 9 |
| `ddim_trailing` | `bindings.DDIM_TRAILING_SAMPLE_METHOD` | 10 |
| `tcd` | `bindings.TCD_SAMPLE_METHOD` | 11 |
| `res_multistep` | `bindings.RES_MULTISTEP_SAMPLE_METHOD` | 12 |
| `res_2s` | `bindings.RES_2S_SAMPLE_METHOD` | 13 |

### 调度器（Scheduler）

| 字符串名称 | Go 常量 | 值 |
|-----------|---------|----|
| `discrete` | `bindings.DISCRETE_SCHEDULER` | 0 |
| `karras` | `bindings.KARRAS_SCHEDULER` | 1 |
| `exponential` | `bindings.EXPONENTIAL_SCHEDULER` | 2 |
| `ays` | `bindings.AYS_SCHEDULER` | 3 |
| `gits` | `bindings.GITS_SCHEDULER` | 4 |
| `sgm_uniform` | `bindings.SGM_UNIFORM_SCHEDULER` | 5 |
| `simple` | `bindings.SIMPLE_SCHEDULER` | 6 |
| `smoothstep` | `bindings.SMOOTHSTEP_SCHEDULER` | 7 |
| `kl_optimal` | `bindings.KL_OPTIMAL_SCHEDULER` | 8 |
| `lcm` | `bindings.LCM_SCHEDULER` | 9 |
| `bong_tangent` | `bindings.BONG_TANGENT_SCHEDULER` | 10 |

### 解析与校验

```go
// 字符串 → 枚举（推荐用于 HTTP 参数、配置文件等外部输入）
method, err := stablediffusion.ParseSampleMethod("euler_a")
sched, err := stablediffusion.ParseScheduler("karras")

// 直接校验枚举值
err := bindings.DPMPP2M_SAMPLE_METHOD.Validate()

// 查询当前加载的库实际支持的支持项
samplers  := stablediffusion.SupportedSampleMethods()
schedulers := stablediffusion.SupportedSchedulers()
```

**错误行为：**

- **传入未知值**（如 `ParseSampleMethod("foo")` 或 `SampleMethod(99)`）：
  返回 `*bindings.UnknownSampleMethodError` / `*bindings.UnknownSchedulerError`，
  错误信息中列出全部支持的取值。
- **库版本过旧**：如果某个采样器/调度器在 Go 绑定中已知，但当前加载的
  `libstable-diffusion` 动态库版本过旧、尚不支持它，返回
  `*bindings.UnsupportedSampleMethodError` / `*bindings.UnsupportedSchedulerError`，
  **错误信息中会明确提示需要升级库文件**
  （`please upgrade the libstable-diffusion library file`）。
  此时重新编译新版 `stable-diffusion.cpp` 动态库即可。
- `GenerateImage` / `GenerateVideo` 在调用底层库之前会自动校验
  `Sampler`（含 `HighNoiseSampler`）配置，非法值直接返回上述错误。

---

## 🛠 目录结构

```
├── bindings/                     # 底层 purego 绑定（1:1 映射 C API）
│   ├── stablediffusion.go        # 绑定实现 + Mock
│   ├── sampler.go                # 采样器/调度器枚举、解析与校验
│   ├── stablediffusion_test.go   # 字符串/回调/Mock 测试
│   └── sampler_test.go           # 采样器/调度器枚举与错误测试
├── stablediffusion.go            # 高层 Go API 封装
├── test/
│   └── stablediffusion_test.go   # 高层 API 集成测试
├── examples/
│   └── server/server.go          # 生产级 HTTP 服务示例
├── stable-diffusion.cpp/         # 核心 C++ 库（子模块）
├── Dockerfile                    # 多阶段 Linux Docker 构建
└── docker-compose.yml            # 一键部署配置
```

---

## 🧪 测试用例详情

### 底层绑定层 (`bindings/`)

| # | 测试名称 | 验证内容 | 预期结果 |
|---|----------|----------|----------|
| 1.1 | `TestCStringGoString` | Go↔C 字符串转换含 NULL 终止符 | 字符串相等，NULL 存在 |
| 1.2 | `TestMockImplementation` | 无动态库时自动降级到 Mock | `GetSystemInfo()`/`GetVersion()` 返回非空 |
| 1.3 | `TestCallbackWrapper` | `purego.NewCallback` 包装不被 GC 回收 | `currentLogCallback != 0` |
| 1.4 | `TestSampleMethodEnumValues` / `TestSchedulerEnumValues` | 枚举值与 `stable-diffusion.h` 一致 | `EULER=0` … `SCHEDULER_COUNT=11` |
| 1.5 | `TestEnumStringNames` / `TestParseRoundTrip` | 名称与 C 库一致、解析往返 | `euler_a`/`karras` 等往返相等 |
| 1.6 | `TestParseUnknownValues` / `TestValidateOutOfRange` | 未知值返回错误 | 错误信息列出支持项 |
| 1.7 | `TestLibraryTooOld` | 模拟旧版动态库 | 错误信息提示升级库文件 |
| 1.8 | `TestSupportedLists` / `TestMockNameConversion` | 支持项查询与 Mock 行为 | 返回全部枚举值 |

### 高层 API 层 (`test/`)

| # | 测试名称 | 验证内容 | 预期结果 |
|---|----------|----------|----------|
| 2.1 | `TestSystemInfo` | `stablediffusion.GetSystemInfo()` | 返回非空字符串 |
| 2.2 | `TestVersionInfo` | `GetVersion()` + `GetCommit()` | 均返回非空 |
| 2.3 | `TestCreateContext` | `NewContext` 无模型文件时 | 返回 `error` |
| 2.4 | `TestCreateUpscaler` | `NewUpscaler` 无模型文件时 | 返回 `error` |
| 2.5 | `TestConvertModel` | `ConvertModel` 无文件时 | 返回 `error` |
| 2.6 | `TestImageGenerationConfig` | `GenerationConfig` 结构体初始化 | `Width=512`, `Height=512` |
| 2.7 | `TestParseSamplerScheduler` | `ParseSampleMethod`/`ParseScheduler` | 合法名称解析正确 |
| 2.8 | `TestParseUnknownSamplerScheduler` | 未知采样器/调度器 | 返回错误并列出支持项 |
| 2.9 | `TestGenerateImageValidatesSampler` | `GenerateImage` 传入非法枚举 | 调用底层库前返回错误 |
| 2.10 | `TestSupportedSamplersSchedulers` | 支持项查询 | 返回全部枚举值 |

### 运行所有测试

```bash
go test -v ./...
```

预期输出：
```
=== RUN   TestCStringGoString
--- PASS: TestCStringGoString
=== RUN   TestMockImplementation
--- PASS: TestMockImplementation
=== RUN   TestCallbackWrapper
--- PASS: TestCallbackWrapper
=== RUN   TestSampleMethodEnumValues
--- PASS: TestSampleMethodEnumValues
=== RUN   TestSchedulerEnumValues
--- PASS: TestSchedulerEnumValues
...
PASS
ok  github.com/example/stablediffusion/bindings

=== RUN   TestSystemInfo
--- PASS: TestSystemInfo
=== RUN   TestVersionInfo
--- PASS: TestVersionInfo
=== RUN   TestCreateContext
--- PASS: TestCreateContext
=== RUN   TestCreateUpscaler
--- PASS: TestCreateUpscaler
=== RUN   TestConvertModel
--- PASS: TestConvertModel
=== RUN   TestImageGenerationConfig
--- PASS: TestImageGenerationConfig
=== RUN   TestParseSamplerScheduler
--- PASS: TestParseSamplerScheduler
=== RUN   TestParseUnknownSamplerScheduler
--- PASS: TestParseUnknownSamplerScheduler
=== RUN   TestGenerateImageValidatesSampler
--- PASS: TestGenerateImageValidatesSampler
=== RUN   TestSupportedSamplersSchedulers
--- PASS: TestSupportedSamplersSchedulers
PASS
ok  github.com/example/stablediffusion/test
```

---

## 🔧 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SD_LIB_PATH` | 动态库完整路径 | 自动搜索平台默认名 |
| `MODEL_PATH` | 模型文件路径 | `models/model.gguf` |
| `PORT` | HTTP 服务端口 | `8080` |
| `HOST` | HTTP 监听地址 | `0.0.0.0` |
| `GENERATE_TIMEOUT` | 生成超时时间 | `5m` |
