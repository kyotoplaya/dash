package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type hardwareConfig struct{}

type hardwareModule struct {
	cfg hardwareConfig
}

func NewHardwareModule() Module {
	return &hardwareModule{
		cfg: hardwareConfig{},
	}
}

func (m *hardwareModule) Config() any {
	return &m.cfg
}

func getCpuTemp() (string, error) {
	psCommand := `
	$zone = Get-CimInstance -Namespace "root/wmi" -ClassName MSAcpi_ThermalZoneTemperature | Select-Object -First 1;
	if ($zone) {
		$celsius = ($zone.CurrentTemperature / 10) - 273.15;
		[math]::Round($celsius, 1)
	}`

	psCommand = strings.ReplaceAll(psCommand, "\n", " ")
	psCommand = strings.ReplaceAll(psCommand, "\t", " ")

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	tempValue := strings.TrimSpace(out.String())
	if tempValue == "" {
		return "", fmt.Errorf("система не вернула данные (датчик заблокирован BIOS/драйвером)")
	}

	return tempValue, nil
}

func getCpuUsage() (string, error) {
	psCommand := `(Get-CimInstance -ClassName Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	usageValue := strings.TrimSpace(out.String())
	if usageValue == "" {
		return "", fmt.Errorf("не удалось получить данные о загрузке процессора")
	}

	return usageValue, nil
}

func getGpuTemp() (string, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader,nounits")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("утилита nvidia-smi не найдена или видеокарта не поддерживается")
	}

	tempValue := strings.TrimSpace(out.String())
	if tempValue == "" {
		return "", fmt.Errorf("не удалось получить данные от видеокарты")
	}

	return tempValue, nil
}

func getGPUCoreLoad() (string, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("не удалось получить загрузку GPU через nvidia-smi")
	}

	usageValue := strings.TrimSpace(out.String())
	if usageValue == "" {
		return "", fmt.Errorf("nvidia-smi вернул пустое значение загрузки GPU")
	}

	// nvidia-smi может вернуть несколько строк для нескольких GPU — берём первую.
	if idx := strings.IndexByte(usageValue, '\n'); idx != -1 {
		usageValue = strings.TrimSpace(usageValue[:idx])
	}

	return usageValue, nil
}

func getGPUVramUsage() (string, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=memory.used", "--format=csv,noheader,nounits")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("не удалось получить VRAM через nvidia-smi")
	}

	vramValue := strings.TrimSpace(out.String())
	if vramValue == "" {
		return "", fmt.Errorf("nvidia-smi вернул пустое значение VRAM")
	}

	if idx := strings.IndexByte(vramValue, '\n'); idx != -1 {
		vramValue = strings.TrimSpace(vramValue[:idx])
	}

	return vramValue + " MB", nil
}

func getRAMMetrics() (float64, float64, error) {
	psCommand := `
	$os = Get-CimInstance -ClassName Win32_OperatingSystem;
	$totalGB = [math]::Round($os.TotalVisibleMemorySize / 1MB, 1);
	$freeGB =  [math]::Round($os.FreePhysicalMemory / 1MB, 1);
	$usedGB =  [math]::Round($totalGB - $freeGB, 1);
	"$usedGB,$totalGB"`

	psCommand = strings.ReplaceAll(psCommand, "\n", " ")
	psCommand = strings.ReplaceAll(psCommand, "\t", " ")

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return 0, 0, err
	}

	output := strings.TrimSpace(out.String())
	parts := strings.Split(output, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректный ответ от системы: %s", output)
	}

	used, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга занятой памяти: %v", err)
	}

	total, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга общей памяти: %v", err)
	}

	return used, total, nil
}

func (*hardwareModule) Id() string   { return "hardware" }
func (*hardwareModule) Name() string { return "Hardware" }
func (m *hardwareModule) Collect() (string, error) {
	cpuTemp, err := getCpuTemp()
	if err != nil {
		cpuTemp = "N/A"
	}

	cpuUsage, err := getCpuUsage()
	if err != nil {
		cpuUsage = "-"
	}

	gpuTemp, err := getGpuTemp()
	if err != nil {
		gpuTemp = "N/A"
	}

	gpuCoreLoad, err := getGPUCoreLoad()
	if err != nil {
		gpuCoreLoad = "-"
	}

	gpuVramUsage, err := getGPUVramUsage()
	if err != nil {
		gpuVramUsage = "N/A"
	}

	usedRAM, totalRAM, err := getRAMMetrics()
	if err != nil {
		usedRAM = 0
		totalRAM = 0
	}

	hwStr := fmt.Sprintf("CPU Temp: %s°C\nCPU Usage: %s%%\nGPU Temp: %s\nGPU Load: %s%%\nVRAM Usage: %s\nRAM Usage: %.1f/%.1f GB", cpuTemp, cpuUsage, gpuTemp, gpuCoreLoad, gpuVramUsage, usedRAM, totalRAM)
	return hwStr, nil
}
func (hardwareModule) Interval() time.Duration {
	return 5 * time.Second
}
