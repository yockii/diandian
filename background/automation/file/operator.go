package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"diandian/background/automation/core"
)

// Operator 文件操作实现
type Operator struct{}

// NewOperator 创建文件操作实例
func NewOperator() *Operator {
	return &Operator{}
}

// CreateFile 创建文件
func (o *Operator) CreateFile(path string, content []byte) *core.OperationResult {
	start := time.Now()
	
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("创建目录失败: %s", dir),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	// 创建文件
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("创建文件失败: %s", path),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("成功创建文件: %s", path),
		map[string]interface{}{
			"path": path,
			"size": len(content),
		},
	)
	result.SetDuration(start)
	return result
}

// CreateDir 创建目录
func (o *Operator) CreateDir(path string) *core.OperationResult {
	start := time.Now()
	
	err := os.MkdirAll(path, 0755)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("创建目录失败: %s", path),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("成功创建目录: %s", path),
		map[string]interface{}{
			"path": path,
		},
	)
	result.SetDuration(start)
	return result
}

// MoveFile 移动文件
func (o *Operator) MoveFile(src, dst string) *core.OperationResult {
	start := time.Now()
	
	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("创建目标目录失败: %s", dstDir),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	// 移动文件
	err := os.Rename(src, dst)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("移动文件失败: %s -> %s", src, dst),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("成功移动文件: %s -> %s", src, dst),
		map[string]interface{}{
			"src": src,
			"dst": dst,
		},
	)
	result.SetDuration(start)
	return result
}

// CopyFile 复制文件
func (o *Operator) CopyFile(src, dst string) *core.OperationResult {
	start := time.Now()
	
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("打开源文件失败: %s", src),
			err,
		)
		result.SetDuration(start)
		return result
	}
	defer srcFile.Close()
	
	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("创建目标目录失败: %s", dstDir),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("创建目标文件失败: %s", dst),
			err,
		)
		result.SetDuration(start)
		return result
	}
	defer dstFile.Close()
	
	// 复制文件内容
	bytesWritten, err := io.Copy(dstFile, srcFile)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("复制文件内容失败: %s -> %s", src, dst),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("成功复制文件: %s -> %s", src, dst),
		map[string]interface{}{
			"src":   src,
			"dst":   dst,
			"bytes": bytesWritten,
		},
	)
	result.SetDuration(start)
	return result
}

// DeleteFile 删除文件
func (o *Operator) DeleteFile(path string) *core.OperationResult {
	start := time.Now()
	
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result := core.NewErrorResult(
			fmt.Sprintf("文件不存在: %s", path),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	// 删除文件
	err := os.Remove(path)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("删除文件失败: %s", path),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("成功删除文件: %s", path),
		map[string]interface{}{
			"path": path,
		},
	)
	result.SetDuration(start)
	return result
}

// DeleteDir 删除目录
func (o *Operator) DeleteDir(path string) *core.OperationResult {
	start := time.Now()
	
	// 检查目录是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result := core.NewErrorResult(
			fmt.Sprintf("目录不存在: %s", path),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	// 删除目录及其内容
	err := os.RemoveAll(path)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("删除目录失败: %s", path),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("成功删除目录: %s", path),
		map[string]interface{}{
			"path": path,
		},
	)
	result.SetDuration(start)
	return result
}

// RenameFile 重命名文件
func (o *Operator) RenameFile(oldPath, newPath string) *core.OperationResult {
	start := time.Now()
	
	// 检查源文件是否存在
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		result := core.NewErrorResult(
			fmt.Sprintf("源文件不存在: %s", oldPath),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	// 重命名文件
	err := os.Rename(oldPath, newPath)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("重命名文件失败: %s -> %s", oldPath, newPath),
			err,
		)
		result.SetDuration(start)
		return result
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("成功重命名文件: %s -> %s", oldPath, newPath),
		map[string]interface{}{
			"old_path": oldPath,
			"new_path": newPath,
		},
	)
	result.SetDuration(start)
	return result
}

// FileExists 检查文件是否存在
func (o *Operator) FileExists(path string) (bool, *core.OperationResult) {
	start := time.Now()
	
	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)
	
	result := core.NewSuccessResult(
		fmt.Sprintf("检查文件存在性: %s (存在: %t)", path, exists),
		map[string]interface{}{
			"path":   path,
			"exists": exists,
		},
	)
	result.SetDuration(start)
	return exists, result
}

// GetFileInfo 获取文件信息
func (o *Operator) GetFileInfo(path string) (interface{}, *core.OperationResult) {
	start := time.Now()
	
	info, err := os.Stat(path)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("获取文件信息失败: %s", path),
			err,
		)
		result.SetDuration(start)
		return nil, result
	}
	
	fileInfo := map[string]interface{}{
		"name":     info.Name(),
		"size":     info.Size(),
		"mode":     info.Mode().String(),
		"mod_time": info.ModTime(),
		"is_dir":   info.IsDir(),
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("获取文件信息: %s", path),
		fileInfo,
	)
	result.SetDuration(start)
	return fileInfo, result
}

// ListDir 列出目录内容
func (o *Operator) ListDir(path string) ([]string, *core.OperationResult) {
	start := time.Now()
	
	entries, err := os.ReadDir(path)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("读取目录失败: %s", path),
			err,
		)
		result.SetDuration(start)
		return nil, result
	}
	
	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("列出目录内容: %s (%d 个项目)", path, len(files)),
		map[string]interface{}{
			"path":  path,
			"count": len(files),
			"files": files,
		},
	)
	result.SetDuration(start)
	return files, result
}

// WriteTextFile 写入文本文件
func (o *Operator) WriteTextFile(path, content string) *core.OperationResult {
	return o.CreateFile(path, []byte(content))
}

// ReadTextFile 读取文本文件
func (o *Operator) ReadTextFile(path string) (string, *core.OperationResult) {
	start := time.Now()
	
	content, err := os.ReadFile(path)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("读取文件失败: %s", path),
			err,
		)
		result.SetDuration(start)
		return "", result
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("读取文件: %s (%d 字节)", path, len(content)),
		map[string]interface{}{
			"path": path,
			"size": len(content),
		},
	)
	result.SetDuration(start)
	return string(content), result
}
