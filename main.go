package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite" // Pure-Go SQLite driver

	"github.com/edhiraj/skt-morph-tool/engine"
)

// Query parameter structs (matches Rust's serde #[derive(Deserialize)])
type VerbQuery struct {
	Root       string `form:"root"`
	Upasarga   string `form:"upasarga"`
	Lakara     string `form:"lakara"`
	Purusha    string `form:"purusha"`
	Voice      string `form:"voice"`
	Prayoga    string `form:"prayoga"`
	Derivative string `form:"derivative"`
}

type ParticipleQuery struct {
	Root       string `form:"root"`
	Upasarga   string `form:"upasarga"`
	Pratyaya   string `form:"pratyaya"`
	Gender     string `form:"gender"`
	Derivative string `form:"derivative"`
}

type DeclensionQuery struct {
	Base   string `form:"base"`
	Gender string `form:"gender"`
}

func main() {
	// 1. Initialize SQLite Database Pool
	// Go's *sql.DB is natively a thread-safe connection pool!
	db, err := sql.Open("sqlite", "data/skt_morphology.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 2. Initialize Gin Router
	r := gin.Default()

	// 3. API Routes
	r.GET("/api/analyze/:word", func(c *gin.Context) {
		word := c.Param("word")
		result := engine.Analyze(db, word)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/api/dhatus/:query", func(c *gin.Context) {
		query := c.Param("query")
		result := engine.SearchDhatu(db, query)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/api/generate/verb", func(c *gin.Context) {
		var q VerbQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := engine.GenerateVerb(db, q.Root, q.Upasarga, q.Lakara, q.Purusha, q.Voice, q.Prayoga, q.Derivative)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/api/generate/participle", func(c *gin.Context) {
		var q ParticipleQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := engine.GenerateParticiple(db, q.Root, q.Upasarga, q.Pratyaya, q.Gender, q.Derivative)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/api/generate/declension", func(c *gin.Context) {
		var q DeclensionQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := engine.GenerateDeclension(q.Base, q.Gender)
		c.JSON(http.StatusOK, result)
	})

	// 4. Serve Static Frontend Files (Axum ServeDir equivalent)
	r.Static("/assets", "frontend/dist/assets")
	r.StaticFile("/", "frontend/dist/index.html")
	
	// Fallback for single-page applications (SPA routing)
	r.NoRoute(func(c *gin.Context) {
		c.File("frontend/dist/index.html")
	})

	// 5. Environment PORT handling
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("✅ Server running on http://%s\n", addr)

	// 6. Start Server
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}