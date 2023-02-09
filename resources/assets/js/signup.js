function submitSignupForm(event) {
    let email = document.getElementById('email').value
    let password = document.getElementById('password').value
    let confirmPassword = document.getElementById('confirm-password').value

    let data = {
        email: email,
        password: password,
    }

    if (password !== confirmPassword) {
        alert('Passwords do not match!')
        return
    }   

    fetch('/user/signup', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-type': 'application/json'
        },
    }).then(res => { 
        if (!res.ok) {
            throw new Error(`HTTP error, status = ${res.status}`)
        }
        alert('Your account was created successfully!')
        window.location.href = '/user/login'
    })
    .catch(err => {
        alert('Something went wrong. Try again.')
    })

    event.preventDefault()
}       