package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"github.com/edhiraj/skt-morph-tool/engine"
)

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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	db, err := sql.Open("sqlite", "data/skt_morphology.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/api/analyze/:word", func(c *gin.Context) {
		word := strings.TrimSpace(c.Param("word"))
		if word == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Word parameter cannot be empty"})
			return
		}
		result := engine.Analyze(db, word)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/api/dhatus/:query", func(c *gin.Context) {
		query := strings.TrimSpace(c.Param("query"))
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

	r.Static("/assets", "frontend/dist/assets")
	r.StaticFile("/", "frontend/dist/index.html")

	r.NoRoute(func(c *gin.Context) {
		c.File("frontend/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("✅ Server running on http://%s\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}