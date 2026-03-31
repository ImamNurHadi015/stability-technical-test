package store

import "stability-test-task-api/models"

var Tasks = []models.Task{
	{ID: 1, Title: "Learn Go", Done: false},
	{ID: 2, Title: "Build API", Done: false},
}

func GetAllTasks() []models.Task {
	return Tasks
}

// Merubah penggunaan pointer ke local variable menjadi pointer ke element asli 
func GetTaskByID(id int) *models.Task {
	for i := range Tasks {
		if Tasks[i].ID == id {
			return &Tasks[i]
		}
	}
	return nil
}


// Mengubah parameter dari value menjadi pointer 
func AddTask(task *models.Task) {
	// Memperbaiki bug auto increment id yang menyebabkan return 0
	maxID := 0
	for i := range Tasks {
		if Tasks[i].ID > maxID {
			maxID = Tasks[i].ID
		}
	}
// Memperbaiki logic id increment
	task.ID = maxID + 1
	Tasks = append(Tasks, *task)
}

// Menambahkan endpoint Update untuk memodifikasi data task
func UpdateTask(id int, updated models.Task) *models.Task {
    for i := range Tasks {
        if Tasks[i].ID == id {
            if updated.Title != "" {
                Tasks[i].Title = updated.Title
            }
            Tasks[i].Done = updated.Done
            return &Tasks[i]
        }
    }
    return nil
}

func DeleteTask(id int) {
	for i, t := range Tasks {
		if t.ID == id {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
			// Menambahkan return untuk mencegah loop tetap berjalan setelah delete yang bisa menyebabkan index out of range karena panjang slice sudah berubah
			return
		}
	}
}
