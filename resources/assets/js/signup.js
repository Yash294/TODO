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

function submitSignupForm(event, form) {
    event.preventDefault()

    if (!form.checkValidity()) {
        return
    }
    let email = document.getElementById('email').value
    let password = document.getElementById('password').value
    let confirmPassword = document.getElementById('confirm-password').value
    const confirmation = document.getElementById('server-side-validation')

    let data = {
        email: email,
        password: password,
    }

    if (password !== confirmPassword) {
        confirmation.className = 'alert alert-danger'
        confirmation.textContent = 'Passwords do not match.'
    }   

    fetch('/user/signup', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Content-type': 'application/json'
        },
    })
    .then(res => res.json())
    .then(result => {
        
        result.message = result.message.charAt(0).toUpperCase() + result.message.slice(1) + '.'
        confirmation.textContent = result.message

        if (!result.success) {
            confirmation.className = 'alert alert-danger'
        } else {
            confirmation.className = 'alert alert-success'
            setTimeout(() => {
                window.location.href = '/user/login'
            }, 2000)
        }
    })
}       