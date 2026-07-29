import api from './request'

export function getNotes(params) {
  return api.get('/notes', { params })
}

export function getNote(id) {
  return api.get(`/notes/${id}`)
}

export function createNote(data) {
  return api.post('/notes', data)
}

export function updateNote(id, data) {
  return api.put(`/notes/${id}`, data)
}

export function deleteNote(id) {
  return api.delete(`/notes/${id}`)
}

export function getTags() {
  return api.get('/notes/tags')
}

export function renameTag(oldName, newName) {
  return api.put('/notes/tags', {
    old_name: oldName,
    new_name: newName
  })
}

export function deleteTag(name) {
  return api.delete('/notes/tags', {
    data: { name }
  })
}

export function uploadImage(file) {
  const formData = new FormData()
  formData.append('file', file)
  return api.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
