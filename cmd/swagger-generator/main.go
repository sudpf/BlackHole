package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/swaggo/swag/gen"
)

func main() {
	var sourceFile string
	var searchDirs string
	var outputDir string

	// Parse command-line arguments
	flag.StringVar(&sourceFile, "source", "", "Main API source file")
	flag.StringVar(&searchDirs, "search", "", "Comma-separated directories to scan")
	flag.StringVar(&outputDir, "output", "docs", "Directory to save the generated Swagger documentation")
	flag.Parse()

	if sourceFile == "" {
		fmt.Println("Source file is required")
		os.Exit(1)
	}

	// Call the function to generate Swagger documentation
	err := generateSwagger(sourceFile, searchDirs, outputDir)
	if err != nil {
		fmt.Printf("Error generating Swagger documentation: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Swagger documentation generated successfully.")
}

// generateSwagger generates Swagger documentation using swag's internal API
func generateSwagger(sourceFile, searchDirs, outputDir string) error {
	dir := filepath.Dir(sourceFile)
	file := filepath.Base(sourceFile)
	if searchDirs == "" {
		searchDirs = dir
	}

	// Create a new generator instance
	generator := gen.New()

	// Set up configuration
	conf := &gen.Config{
		// 遍历需要查询注释的目录
		SearchDir: searchDirs,
		// 不包含哪些文件
		Excludes: "",
		// 输出目录
		OutputDir:   outputDir,
		OutputTypes: []string{"go", "json", "yaml"},
		// 整个swagger接口的说明文档注释
		MainAPIFile: file,
		// 名字的显示策略，比如首字母大写等
		PropNamingStrategy: "",
		// 是否要解析vendor目录
		ParseVendor: false,
		// 是否要解析外部依赖库的包
		ParseDependency: false,
		// 是否要解析标准库的包
		ParseInternal: true,
		// 是否要查找markdown文件，这个markdown文件能用来为tag增加说明格式
		MarkdownFilesDir: "",
		// 是否应该在docs.go中生成时间戳
		GeneratedTime: false,
		ParseDepth:    2,
		Strict:        true,
	}

	// Run the generator
	err := generator.Build(conf)
	if err != nil {
		return fmt.Errorf("failed to build Swagger documentation: %w", err)
	}

	return nil
}
