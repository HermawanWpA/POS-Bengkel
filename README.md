# POS Echo Application - Backend API

Repositori ini berisi kode sumber untuk layanan backend API **Point of Sale (POS) & Manajemen Pelanggan** menggunakan bahasa pemrograman **Go (Golang)**. Service ini dibangun dengan framework **Echo**, menggunakan **GORM** sebagai ORM untuk komunikasi ke database MySQL, dan menerapkan arsitektur bersih (*Clean Architecture*).

---

## 🚀 Fitur Utama
* **Manajemen Pelanggan & Kendaraan**: Relasi *one-to-many* antara data pelanggan dan daftar kendaraan dengan aksi cascade.
* **Server Berkinerja Tinggi**: Menggunakan engine framework **Echo v4** yang efisien dan cepat.
* **Paginasi & Pencarian Dinamis**: Mendukung data server-side pagination dan keyword searching terintegrasi untuk integrasi frontend UI.
* **RESTful API Standards**: Output response terstruktur dan penanganan error terpusat (*Global Error Handling*).

---

## 🛠️ Prasyarat (Prerequisites)
Sebelum menjalankan aplikasi, pastikan komputer Anda telah terinstal modul-modul berikut:
* **Go Compiler** (Versi 1.20 atau yang lebih baru)
* **Laragon** / **XAMPP** (Untuk mengaktifkan database MySQL/MariaDB)
* **Git** (Opsional, untuk manajemen repositori)

---

## 📦 Panduan Instalasi & Setup

### 1. Kloning Repositori (Jika berlaku)
Jika proyek Anda berada dalam git, buka terminal (Git Bash/CMD) lalu arahkan ke direktori Laragon www Anda:
```bash
cd C:/laragon/www
# Atau silakan masuk langsung ke folder proyek yang sudah ada:
cd pos-echo-app
