// Example starter JavaScript for disabling form submissions if there are invalid fields
(function() {
    'use strict';
    window.addEventListener('load', function() {
      // Fetch all the forms we want to apply custom Bootstrap validation styles to
      let forms = document.getElementsByClassName('needs-validation');
      // Loop over them and prevent submission
      Array.prototype.filter.call(forms, function(form) {
        form.addEventListener('submit', function(event) {
          if (form.checkValidity() === false) {
            event.preventDefault();
            event.stopPropagation();
          }
          form.classList.add('was-validated');
        }, false);
      });
    }, false);
})();

let oldTaskName, deleteTaskName

function setEditValues(taskName, description, isDone) {
    oldTaskName = taskName 

    document.getElementById('edit-task-name').value = taskName
    document.getElementById('edit-description').value = description
    document.getElementById('edit-is-done').checked = isDone
}

function setDeleteValue(taskName) {
    deleteTaskName = taskName
}

function submitAddForm(event, form) {
    event.preventDefault()

    if (!form.checkValidity()) {
        const confirmation = document.getElementById('server-side-validation')
        confirmation.className = 'alert alert-danger'
        confirmation.textContent = 'Please specify a task name!'
        return
    }
    let taskName = document.getElementById('add-task-name').value
    let description = document.getElementById('add-description').value

    let data = {
        taskName: taskName,
        description: description
    }

    fetch('/task/add', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-type': 'application/json'
        },
    })
    .then(res => res.json())
    .then(result => {
        const confirmation = document.getElementById('server-side-validation')
        result.message = result.message.charAt(0).toUpperCase() + result.message.slice(1) + '.'
        confirmation.textContent = result.message

        if (!result.success) {
            confirmation.className = 'alert alert-danger'
        } else {
            confirmation.className = 'alert alert-success'
            setTimeout(() => {
                confirmation.className = ''
                confirmation.textContent = ''
            }, 5000)

            const task = document.createElement('li')
            task.className = 'list-group-item list-group-item-dark'
            task.id = taskName

            const name = document.createElement('h4')
            name.textContent = taskName

            const taskDescription = document.createElement('p')
            taskDescription.textContent = description

            const taskCompleted = document.createElement('span')
            taskCompleted.className = 'badge badge-danger'
            taskCompleted.textContent = 'Incomplete'

            const deleteButton = document.createElement('button')
            deleteButton.className = 'btn btn-primary float-right'
            deleteButton.type = 'button'
            deleteButton.onclick = function() {
                setDeleteValue(name.textContent)
            }
            deleteButton.setAttribute("data-toggle", "modal")
            deleteButton.setAttribute("data-target", "#delete-task")

            const deleteIcon = document.createElement('i')
            deleteIcon.className = "bi bi-trash"

            const editButton = document.createElement('button')
            editButton.className = 'btn btn-primary float-right mr-2'
            editButton.type = 'button'
            editButton.onclick = function() {
                if (taskCompleted.textContent = 'Incomplete') {
                    setEditValues(name.textContent, taskDescription.textContent, false)
                } else {
                    setEditValues(name.textContent, taskDescription.textContent, true)
                }
            }
            editButton.setAttribute("data-toggle", "modal")
            editButton.setAttribute("data-target", "#edit-task")

            const editIcon = document.createElement('i')
            editIcon.className = "bi bi-pencil"

            // build html tree
            deleteButton.appendChild(deleteIcon)
            editButton.appendChild(editIcon)

            task.appendChild(name)
            task.appendChild(taskDescription)
            task.appendChild(taskCompleted)
            task.appendChild(deleteButton)
            task.appendChild(editButton)

            document.getElementById('taskList').append(task)
        } 
    })
}

function submitEditForm(event, form) {
    event.preventDefault()
    
    if (!form.checkValidity()) {
        const confirmation = document.getElementById('server-side-validation')
        confirmation.className = 'alert alert-danger'
        confirmation.textContent = 'Please specify a task name!'
        return
    }

    let taskName = document.getElementById('edit-task-name').value
    let description = document.getElementById('edit-description').value
    let isDone = document.getElementById('edit-is-done').checked

    let data = {
        oldTaskName: oldTaskName,
        taskName: taskName,
        description: description,
        isDone: isDone
    }

    fetch('/task/edit', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-type': 'application/json'
        },
    })
    .then(res => res.json())
    .then(result => {
        const confirmation = document.getElementById('server-side-validation')
        result.message = result.message.charAt(0).toUpperCase() + result.message.slice(1) + '.'
        confirmation.textContent = result.message

        if (!result.success) {
            confirmation.className = 'alert alert-danger'
        } else {
            confirmation.className = 'alert alert-success'
            setTimeout(() => {
                confirmation.className = ''
                confirmation.textContent = ''
            }, 5000)

            for (let taskHeading of document.querySelectorAll('.list-group .list-group-item > h4')) {
                if (taskHeading.textContent == oldTaskName) {
                    const parent = taskHeading.parentNode
    
                    taskHeading.textContent = taskName
                    parent.querySelectorAll('p')[0].textContent = description
                    const badge = parent.querySelectorAll('span')[0]
                    if (isDone) {
                        badge.textContent = 'Complete'
                        badge.className = 'badge badge-success'
                    } else {
                        badge.textContent = 'Incomplete'
                        badge.className = 'badge badge-danger'
                    }
                }
            }
        }
    })
}

function submitDeleteForm(event) {
    event.preventDefault()

    let data = {
        taskName: deleteTaskName,
    }

    fetch('/task/delete', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-type': 'application/json'
        },
    })
    .then(res => res.json())
    .then(result => {
        const confirmation = document.getElementById('server-side-validation')
        result.message = result.message.charAt(0).toUpperCase() + result.message.slice(1) + '.'
        confirmation.textContent = result.message

        if (!result.success) {
            confirmation.className = 'alert alert-danger'    
        } else {
            confirmation.className = 'alert alert-success'
            setTimeout(() => {
                confirmation.className = ''
                confirmation.textContent = ''
            }, 5000)

            for (let taskHeading of document.querySelectorAll('.list-group .list-group-item > h4')) {
                if (taskHeading.textContent == deleteTaskName) {
                    taskHeading.parentNode.parentNode.removeChild(taskHeading.parentNode)
                }
            }
        }
    })
}

function submitLogout() {
    fetch('/user/logout')
    .then(res => res.json())
    .then(result => {
        if (!result.success) {
            const confirmation = document.getElementById('server-side-validation')
            result.message = result.message.charAt(0).toUpperCase() + result.message.slice(1) + '.'
            confirmation.textContent = result.message
            confirmation.className = 'alert alert-danger'    
        } else {
            window.location.href = '/user/login' 
        }
    })
}