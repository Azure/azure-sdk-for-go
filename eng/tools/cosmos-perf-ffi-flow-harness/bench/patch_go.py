#!/usr/bin/env python3
import pathlib
from patch_native import once


def patch_user_agent_test(sdk):
    path = sdk / "user_agent_native_test.go"
    text = path.read_text()
    text = once(text, """\t\t\tsuffixPattern := ` ` + regexp.QuoteMeta(tt.applicationID)
\t\t\tif tt.applicationID == "" {
\t\t\t\tsuffixPattern = ` rustc/[^ |]+`
\t\t\t}
""", """\t\t\tsuffixPattern := ""
\t\t\tif tt.applicationID != "" {
\t\t\t\tsuffixPattern = ` ` + regexp.QuoteMeta(tt.applicationID)
\t\t\t}
""")
    text = once(text, r"(?:^| )azsdk-rust-cosmos-driver/\S+(?: |$)",
                r" \(drv=1\.1\.0-beta\.3; linux; x86_64; [^; )]+; ft=[A-Za-z0-9_-]+\)")
    text = once(text, r"suffixPattern+`\|F[0-9A-Fa-f]+$`", r"`\)`+suffixPattern+`$`")
    path.write_text(text)


def patch_sdk(sdk, diagnostic):
    path = sdk / "driver.go"
    path.write_text(once(path.read_text(), 'const nativeDriverVersion = "0.2.0"',
                        'const nativeDriverVersion = "0.3.0"'))
    patch_user_agent_test(sdk)
    if not diagnostic:
        return
    inputs = pathlib.Path(__file__).parent
    for source, target in (("flow.go", "experimental_flow.go"), ("flow_cgo.go", "experimental_flow_cgo.go")):
        (sdk / target).write_text((inputs / source).read_text())
    path = sdk / "reactor_native.go"
    text = once(path.read_text(), '\t"sync"\n', '\t"sync"\n\t"time"\n')
    text = once(text, "type pendingOperation struct {\n", "type pendingOperation struct {\n\tflowDeliveredAt time.Time\n")
    text = once(text, "\tp.result <- result\n",
                "\tp.flowDeliveredAt = time.Now()\n\tflowEvent(flowDelivered, 1)\n\tp.result <- result\n")
    text = once(text, "\t\tC.uint32_t(timeoutMillis),\n\t)\n\tif n == 0 {",
                "\t\tC.uint32_t(timeoutMillis),\n\t)\n\tflowWait(int(n))\n\tif n == 0 {")
    path.write_text(text)
    path = sdk / "operation_native.go"
    text = path.read_text()
    start = text.index("func (d *nativeDriver) awaitCompletion(")
    end = text.index("\nfunc awaitOperationResult(", start)
    part = once(text[start:end], "\top := submit(d.reactor.queue, C.intptr_t(handle), &preError)\n",
                "\tflowEvent(flowSubmitAttempts, 1)\n\top := submit(d.reactor.queue, C.intptr_t(handle), &preError)\n")
    part = once(part, "\tif op == nil {\n", "\tif op == nil {\n\t\tflowEvent(flowSubmitErrors, 1)\n")
    part = once(part, "\t}\n\tdefer func() {\n", "\t}\n\tflowEvent(flowSubmitAccepted, 1)\n\tdefer func() {\n")
    text = text[:start] + part + text[end:]
    text = once(text, "\tcase result := <-pending.result:\n\t\treturn result.prepare(), nil",
                "\tcase result := <-pending.result:\n\t\tflowResumed(pending.flowDeliveredAt)\n\t\treturn result.prepare(), nil")
    text = once(text, "\t\tif result, ok := pending.abandon(); ok {\n\t\t\treturn result.prepare(), nil",
                "\t\tif result, ok := pending.abandon(); ok {\n\t\t\tflowResumed(pending.flowDeliveredAt)\n\t\t\treturn result.prepare(), nil")
    path.write_text(text)


def patch_adapter(adapter, diagnostic):
    inputs = pathlib.Path(__file__).parent
    path = adapter / "sdk_v2.go"
    path.write_text(once(
        path.read_text(),
        "\t\t\tRuntime: runtime,\n\t\t\tDiagnosticsVerbosity: azcosmos.DiagnosticsVerbositySummary,",
        "\t\t\tRuntime: runtime,\n\t\t\tDiagnosticsVerbosity: azcosmos.DiagnosticsVerbositySummary,\n"
        "\t\t\tOperation: azcosmos.OperationOptions{BinaryEncoding: &azcosmos.BinaryEncodingOptions{Disabled: true}},"))
    path = adapter / "main.go"
    text = once(path.read_text(), "func runMain(arguments []string) error {",
                "func runMain(arguments []string) (runErr error) {")
    text = once(text, "\tdefer func() { _ = client.Close() }()\n", """\tdefer func() {
        runErr = errors.Join(runErr, client.Close())
        if directory := os.Getenv("COSMOS_PERF_FLOW_DIR"); directory != "" {
            runErr = errors.Join(runErr, experimentFlowWrite(directory))
        }
    }()
""")
    text = once(text, "\t_, err := client.Read(ctx, itemPK(value.owner, 0), itemID(value.owner, 0))\n",
                "\tread, err := client.Read(ctx, itemPK(value.owner, 0), itemID(value.owner, 0))\n")
    text = once(text, "\tprobe, err := signedHTTPProbe(ctx, value)\n",
                "\tif err := validateDocument(read.Body, value.owner, 0); err != nil {\n"
                "\t\treturn nil, fmt.Errorf(\"SDK preflight payload validation: %w\", err)\n\t}\n"
                "\tprobe, err := signedHTTPProbe(ctx, value)\n")
    text = once(text, "func runBenchmark(ctx context.Context, client sdkClient, value options) error {",
                "func runBenchmark(ctx context.Context, client sdkClient, value options) (runErr error) {")
    old = """\tif err := startDiagnostics(); err != nil {
\t\treturn err
\t}
\tdefer stopDiagnostics()

\tif err := runPhase(ctx, client, value, value.warmup, false, nil); err != nil {
\t\treturn fmt.Errorf("warmup: %w", err)
\t}
\tcpuStart := processCPUSeconds()
\tstarted := time.Now()
"""
    new = """\tif err := experimentPhase(1); err != nil { return err }
    if err := runPhase(ctx, client, value, value.warmup, false, nil); err != nil {
        return fmt.Errorf("warmup: %w", err)
    }
    if err := measurementBarrier("start"); err != nil { return err }
    defer func() { runErr = errors.Join(runErr, stopDiagnostics()) }()
    if err := startDiagnostics(); err != nil { return err }
    if err := experimentPhase(2); err != nil { return err }
    cpuStart := processCPUSeconds()
    monotonicStart, err := monotonicNS()
    if err != nil { return err }
    startUnixNS := time.Now().UnixNano()
    started := time.Now()
"""
    text = once(text, old, new)
    text = once(text, "\tcpuSeconds := processCPUSeconds() - cpuStart\n", """    monotonicEnd, err := monotonicNS()
    if err != nil { return err }
    cpuSeconds := processCPUSeconds() - cpuStart
    if err := experimentPhase(0); err != nil { return err }
    if err := stopDiagnostics(); err != nil { return err }
    if err := measurementBounds(monotonicStart, monotonicEnd); err != nil { return err }
    if err := measurementBarrier("end"); err != nil { return err }
""")
    text = once(text, 'map[string]any{"wallSeconds": wallSeconds, "processCPUSeconds": cpuSeconds}',
                'map[string]any{"wallSeconds": wallSeconds, "processCPUSeconds": cpuSeconds, "measurementStartUnixNs": startUnixNS}')
    begin = text.index("\nvar (\n\tcpuProfileFile")
    end = text.index("\nfunc writeJSON(", begin)
    text = text[:begin] + text[end:]
    for entry in ('\t"runtime"\n', '\t"runtime/pprof"\n', '\truntimeTrace "runtime/trace"\n'):
        text = once(text, entry, "")
    path.write_text(text)
    for name in ("barrier_go.go", "diagnostics_go.go"):
        (adapter / name).write_text((inputs / name).read_text())
    if diagnostic:
        hooks = """package main
import "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
func experimentPhase(phase uint32) error { return azcosmos.ExperimentalFlowPhase(phase) }
func experimentFlowWrite(directory string) error { return azcosmos.ExperimentalFlowWrite(directory) }
"""
    else:
        hooks = """package main
func experimentPhase(uint32) error { return nil }
func experimentFlowWrite(string) error { return nil }
"""
    (adapter / "experimental_hooks.go").write_text(
        "// Copyright (c) Microsoft Corporation. All rights reserved.\n"
        "// Licensed under the MIT License.\n\n" + hooks)
