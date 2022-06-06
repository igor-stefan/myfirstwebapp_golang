function bootstrap_form_validation() {
    'use strict'
    // Fetch all the forms we want to apply custom Bootstrap validation styles to
    let forms = document.querySelectorAll('.needs-validation');

    // Loop over them and prevent submission
    Array.prototype.slice.call(forms)
        .forEach(function (form) {
            form.addEventListener('submit', function (event) {
                if (!form.checkValidity()) {
                    event.preventDefault();
                    event.stopPropagation();
                }
                form.classList.add('was-validated');
            }, false)
        })
}

const alertar = myNotifications();
function myNotifications() { //cria um modulo myNotifications() com vários tipos de notificação
    const Toast = function (c) { //Toast -> rápido aviso no canto da tela
        const {
            msg = "",
            icon = "success",
            position = "bottom-end",
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 2500,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })
        Toast.fire({});
    }

    const Ok = function (c) { // Ok -> Animacao com marcacao do "vêzinho" de ok, verificado
        const {
            msg = "Ok",
            icon = "success",
            text = "Tudo certo!",
            footer = "",
        } = c;
        Swal.fire({
            icon: icon,
            title: msg,
            text: text,
            footer: footer,
        })
    }

    const erro = function (c) { // erro -> notificacao de erro com um "xis vermelho e uma msg" na tela
        const {
            msg = "Erro!",
            icon = "error",
            text = "Algo deu errado...",
            footer = "",
        } = c;
        Swal.fire({
            icon: icon,
            title: msg,
            text: text,
            footer: footer,
        })
    }

    async function custom(c) { // custom -> abre um quadro com um html e msgs diversas
        const {
            icon = "",
            msg = "vazio...",
            title = "nada foi declarado",
            showConfirmButton = true,
        } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) { c.willOpen() }
            },
            preConfirm: () => {
                if (c.preConfirm !== undefined) { c.preConfirm() }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) { c.didOpen() }
            },
        })

        if (result) { //há um resultado
            if (result.dismiss !== Swal.DismissReason.cancel) { //nao foi do botao cancelar
                if (result.value !== "") {//o valor é diferente de vazio
                    if (c.callback !== undefined) { //há um callback
                        c.callback(result); //chama o callback com o resultado
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }
    return { // tipos de chamada da funcao alertar
        toast: Toast,
        sucesso: Ok,
        erro: erro,
        custom: custom,
    }
}

function notify(msg, msgType) {
    notie.alert({
        type: msgType,
        text: msg
    })
}