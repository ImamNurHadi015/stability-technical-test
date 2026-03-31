# Stability Task Manager API

Submission technical test untuk Fullstack Developer Intern di  PT. Tirtamas Coldstorindo Logistik

## Tech Stack

- **Go** — Bahasa pemrograman backend
- **Fiber** — Web framework
- **In-memory storage** — Penyimpanan data sederhana

## Cara Menjalankan

```bash
# Clone repository
git clone https://github.com/ImamNurHadi015/stability-technical-test.git
cd stability-technical-test

# Jalankan aplikasi
go run main.go
```

Server akan berjalan di `http://127.0.0.1:3000`

## API Endpoints

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/tasks` | Mengambil semua task |
| GET | `/tasks/:id` | Mengambil task berdasarkan ID |
| POST | `/tasks` | Membuat task baru |
| DELETE | `/tasks/:id` | Menghapus task berdasarkan ID |

---

## Bug yang Ditemukan & Diperbaiki

### 1. Incorrect Status Code pada `GET /tasks/:id`
**Kategori:** Incorrect status codes

**Masalah:**  
Ketika mengambil task dengan ID yang tidak ada (contoh `/tasks/3`), API mengembalikan status code `200 OK` padahal seharusnya `404 Not Found`.

```go
// Sebelum
if task == nil {
    return c.Status(200).JSON(fiber.Map{
        "error": "task not found",
    })
}
```

**Perbaikan:**  
Mengubah status code menjadi `404` ketika task tidak ditemukan.

```go
// Sesudah
if task == nil {
    return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "error": "task not found",
    })
}
```

---

### 2. Incorrect Status Code & Missing Validation pada `DELETE /tasks/:id`
**Kategori:** Incorrect status codes, Missing validation

**Masalah:**  
Ketika menghapus task dengan ID yang tidak ada (contoh `/tasks/99`), API tetap mengembalikan `200 OK` dengan pesan `"deleted"` meskipun task tidak ditemukan.

**Perbaikan:**  
Menambahkan validasi untuk memeriksa apakah task ada sebelum dihapus. Mengembalikan `404 Not Found` jika task tidak ditemukan.

```go
// Sesudah
task := store.GetTaskByID(id)
if task == nil {
    return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "error": "task not found",
    })
}
store.DeleteTask(id)
```

---

### 3. Endpoint Returning Incorrect Data pada `POST /tasks`
**Kategori:** Endpoint returning incorrect data

**Masalah:**  
Ketika membuat task baru, API selalu mengembalikan `id: 0` dan bukan ID yang seharusnya. Hal ini terjadi karena fungsi `AddTask` tidak pernah melakukan assignment ID pada task, sehingga nilainya tetap zero value Go yaitu `0`.

**Perbaikan:**  
Memperbaiki fungsi `AddTask` agar ID di-generate secara otomatis berdasarkan ID terbesar yang ada. Parameter juga diubah menjadi pointer agar perubahan ID langsung tercermin ke caller.

```go
// Sesudah
func AddTask(task *models.Task) {
    maxID := 0
    for i := range Tasks {
        if Tasks[i].ID > maxID {
            maxID = Tasks[i].ID
        }
    }
    task.ID = maxID + 1
    Tasks = append(Tasks, *task)
}
```

> Menggunakan `maxID + 1` dibanding `len(Tasks) + 1` untuk mencegah duplikasi ID setelah task dihapus.

---

## Improvement

### 1. Input Validation pada `POST /tasks`
**Kategori:** Add input validation, Improve error handling

Menambahkan validasi agar `title` tidak boleh kosong saat membuat task baru. Mengembalikan `400 Bad Request` dengan pesan error yang jelas jika validasi gagal.

```go
// Validasi ditambahkan
if task.Title == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": "title is required",
    })
}
```

Contoh response ketika title kosong:
```json
{
    "error": "title is required"
}
```

### 2. Perbaikan Invalid Pointer pada `GetTaskByID`

Kode sebelumnya mengembalikan `&t` yang merupakan pointer ke salinan lokal dari loop, bukan ke elemen asli di slice. Ini adalah latent bug yang dapat menyebabkan masalah jika ada fitur seperti UPDATE di masa mendatang.

```go
// Sebelum — pointer ke salinan loop
for _, t := range Tasks {
    if t.ID == id {
        return &t
    }
}

// Sesudah — pointer ke elemen asli
for i := range Tasks {
    if Tasks[i].ID == id {
        return &Tasks[i]
    }
}
```

### 3. Menambahkan `return` Setelah Delete pada `DeleteTask`

Menambahkan `return` setelah task dihapus agar loop tidak terus berjalan setelah slice dimodifikasi, yang dapat menyebabkan index out of range.

```go
// Sesudah
func DeleteTask(id int) {
    for i, t := range Tasks {
        if t.ID == id {
            Tasks = append(Tasks[:i], Tasks[i+1:]...)
            return
        }
    }
}
```
