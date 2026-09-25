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

    // 图像生成（采样器和调度器均使用高层枚举常量）
    cfg := stablediffusion.GenerationConfig{
        Prompt:     "A cyberpunk city, 8k wallpaper",
        Width:      512,
        Height:     512,
        BatchCount: 1,
        Sampler: stablediffusion.SamplerConfig{
            Steps:     20,
            TxtCfg:    7.0,
            Method:    stablediffusion.SamplerEulerA,
            Scheduler: stablediffusion.SchedulerKarras,
        },
    }

    images, err := ctx.GenerateImage(cfg)
    if err != nil {
        log.Fatalf("Failed: %v", err)
    }
    log.Printf("Generated %d images", len(images))
}
```

### 支持的采样器（Sampler）与调度器（Scheduler）

采样器与调度器在 `stablediffusion` 包中整理为类型安全的 Go 枚举常量
（底层与 `bindings.SampleMethod` / `bindings.Scheduler` 一一对应，可互换使用）。

**采样器**（`SamplerConfig.Method`，类型 `stablediffusion.Sampler`）：

| Go 常量 | 字符串名称 (`SamplerName`) | 备注 |
|---------|----------------------------|------|
| `SamplerEuler` | `euler` | |
| `SamplerEulerA` | `euler_a` | 最常用的默认采样器之一 |
| `SamplerHeun` | `heun` | |
| `SamplerDPM2` | `dpm2` | |
| `SamplerDPMpp2SA` | `dpmpp_2s_a` | |
| `SamplerDPMpp2M` | `dpmpp_2m` | |
| `SamplerDPMpp2Mv2` | `dpmpp_2m_v2` | |
| `SamplerIPNDM` | `ipndm` | |
| `SamplerIPNDMv` | `ipndm_v` | |
| `SamplerLCM` | `lcm` | LCM 蒸馏模型建议搭配 `SchedulerLcm` |
| `SamplerDDIMTrailing` | `ddim_trailing` | |
| `SamplerTCD` | `tcd` | 需要较新版本的库 |
| `SamplerResMultistep` | `res_multistep` | 需要较新版本的库 |
| `SamplerRes2s` | `res_2s` | 需要较新版本的库 |

**调度器**（`SamplerConfig.Scheduler`，类型 `stablediffusion.Scheduler`）：

| Go 常量 | 字符串名称 (`SchedulerName`) | 备注 |
|---------|------------------------------|------|
| `SchedulerDiscrete` | `discrete` | |
| `SchedulerKarras` | `karras` | 最常用的默认调度器之一 |
| `SchedulerExponential` | `exponential` | |
| `SchedulerAys` | `ays` | |
| `SchedulerGits` | `gits` | |
| `SchedulerSgmUniform` | `sgm_uniform` | |
| `SchedulerSimple` | `simple` | Flow/Flux 系模型常用 |
| `SchedulerSmoothstep` | `smoothstep` | |
| `SchedulerKlOptimal` | `kl_optimal` | |
| `SchedulerLcm` | `lcm` | LCM 蒸馏模型建议搭配 `SamplerLCM` |
| `SchedulerBongTangent` | `bong_tangent` | 需要较新版本的库 |

> 运行时可用 `stablediffusion.Samplers()` / `stablediffusion.Schedulers()` 列出全部
> 已知枚举，用 `IsSamplerSupported` / `IsSchedulerSupported` 查询当前加载的动态库
> 是否支持；也可用 `ParseSampler("euler_a")` / `ParseScheduler("karras")`
> 按名称解析（大小写不敏感）。

**错误处理**：

- 传入未知枚举值（不在上表中的数值）时，`GenerateImage` / `GenerateVideo`
  会在调用动态库**之前**返回错误，例如：
  `unknown sampler value 9999; supported samplers are: euler, euler_a, ...`
- 如果枚举值合法、但当前加载的 `stable-diffusion` 动态库版本过旧而不支持
  （常见于新加的采样器/调度器），错误信息会明确提示**升级库文件**，例如：
  `sampler "res_2s" is not supported by the loaded stable-diffusion library;
  please upgrade the library file (libstable-diffusion.so /
  libstable-diffusion.dylib / stable-diffusion.dll) to a newer version ...`
- 字符串解析同理：未知名称返回 `unknown sampler/scheduler ...` 错误。

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

## 🛠 目录结构

```
├── bindings/                     # 底层 purego 绑定（1:1 映射 C API）
│   ├── stablediffusion.go        # 绑定实现 + Mock
│   ├── enums.go                  # 采样器/调度器枚举、名称与库支持探测
│   └── stablediffusion_test.go   # 字符串/回调/Mock 测试
├── stablediffusion.go            # 高层 Go API 封装
├── enums.go                      # 高层采样器/调度器常量与参数校验
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

### 高层 API 层 (`test/`)

| # | 测试名称 | 验证内容 | 预期结果 |
|---|----------|----------|----------|
| 2.1 | `TestSystemInfo` | `stablediffusion.GetSystemInfo()` | 返回非空字符串 |
| 2.2 | `TestVersionInfo` | `GetVersion()` + `GetCommit()` | 均返回非空 |
| 2.3 | `TestCreateContext` | `NewContext` 无模型文件时 | 返回 `error` |
| 2.4 | `TestCreateUpscaler` | `NewUpscaler` 无模型文件时 | 返回 `error` |
| 2.5 | `TestConvertModel` | `ConvertModel` 无文件时 | 返回 `error` |
| 2.6 | `TestImageGenerationConfig` | `GenerationConfig` 结构体初始化 | `Width=512`, `Height=512` |
| 2.7 | `TestGenerateImageRejectsUnknownSampler` | 未知采样器值 | 返回 `unknown sampler ...` 错误 |
| 2.8 | `TestGenerateImageRejectsUnknownScheduler` | 未知调度器值 | 返回 `unknown scheduler ...` 错误 |
| 2.9 | `TestGenerateVideoRejectsUnknownSampler` | 视频生成传入未知采样器 | 返回 `unknown sampler ...` 错误 |
| 2.10 | `TestParseSamplerAndScheduler` | 名称解析枚举 | 合法名称解析成功，非法名称报错 |

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
