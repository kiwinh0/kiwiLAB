package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Rutas para archivos estáticos
	fs := http.FileServer(http.Dir("../../ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Ruta principal (Dashboard)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<html>
			<head>
				<title>kiwiLAB | Dashboard</title>
				<script src="https://cdn.tailwindcss.com"></script>
			</head>
			<body class="bg-black text-white flex items-center justify-center h-screen">
				<div class="text-center">
					<h1 class="text-5xl font-bold text-red-600">kiwiLAB</h1>
					<p class="mt-4 text-gray-400">Debian 13 x86_64 - Entorno de Desarrollo Listo</p>
				</div>
			</body>
			</html>
		`)
	})

	fmt.Println("🚀 kiwiLAB ejecutándose en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
