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

function submitLoginForm(event, form) {
    if (!form.checkValidity()) {
        return
    }
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
                window.location.href = '/task/dashboard'
            }, 2000)
        }
    })

    event.preventDefault()
}

function submitResetPasswordForm(event, form) {
    if (!form.checkValidity()) {
        return
    }
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
        }
    })

    event.preventDefault()
}