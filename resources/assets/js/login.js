function submitLoginForm(event) {
    let email = document.getElementById('email').value
    let password = document.getElementById('password').value

    let data = {
        email: email,
        password: password,
    }

    fetch('/user/login', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-type': 'application/json'
        },
    }).then(res => {
        if (!res.ok) {
            throw new Error(`HTTP error, status = ${res.status}`)
        }
        alert('You have logged in successfully!')
        window.location.href = '/task/dashboard'
    })
    .catch(err => {
        alert('Something went wrong. Try again.')
    })

    event.preventDefault()
}

function submitResetPasswordForm(event) {
    let resetEmail = document.getElementById('reset-email').value
    let resetPassword = document.getElementById('reset-password').value
    let newPassword = document.getElementById('new-password').value

    let data = {
        email: resetEmail,
        password: resetPassword,
        newPassword: newPassword
    }

    fetch('/user/resetPassword', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-type': 'application/json'
        },
    }).then(res => {
        if (!res.ok) {
            throw new Error(`HTTP error, status = ${res.status}`)
        }
        alert('Your password has been reset successfully!')
        window.location.href = '/user/login'
    })
    .catch(err => {
        alert('Something went wrong. Try again.')
    })

    event.preventDefault()
}