const form = document.querySelector('form')
const input = form.querySelector('input')
const submitButton = form.querySelector('button')
const canvas = document.querySelector('canvas')

form.addEventListener('submit', onSubmit)

/**
 * @param {Event} ev
 */
function onSubmit(ev) {
    ev.preventDefault()

    const data = new FormData(form)

    input.disabled = true
    submitButton.disabled = true

    upload(data).then(statements => {
        renderGraph(statements)
        form.reset()
    }).catch(err => {
        console.error(err)
        alert(err.message)
    }).then(() => {
        input.disabled = false
        submitButton.disabled = false
    })
}

/**
 * @param {FormData} data
 */
function upload(data) {
    return fetch(`/api/upload`, {
        method: 'POST',
        body: data,
    }).then(handleResponse).then(res => res.body)
}

/**
 * @param {Response} res
 */
async function handleResponse(res) {
    const body = await res.clone().json().catch(() => res.text())
    const response = {
        statusCode: res.status,
        statusText: res.statusText,
        headers: res.headers,
        body,
    }

    if (!res.ok) {
        const message = typeof body === 'object' && body !== null && typeof body.message === 'string' && body.message !== ''
            ? body.message
            : typeof body === 'string' && body !== ''
                ? body
                : res.statusText
        throw Object.assign(new Error(message), response)
    }

    return response
}

/**
 * @param {Statement[]} statements
 */
function renderGraph(statements) {
    canvas.innerHTML = ''
    const ctx = canvas.getContext('2d')
    new window['Chart'](ctx, {
        type: 'bar',
        data: {
            labels: ['credit', 'debit'],
            datasets: statements.map(statement => ({
                // label: 'Bank Report',
                data: [
                    statement.credit,
                    statement.debit,
                ],
                backgroundColor: [
                    'rgba(255, 99, 132, 0.6)',
                    'rgba(54, 162, 235, 0.6)',
                ],
                borderWidth: 1,
                borderColor: '#777',
            })),
        },
        options: {
            title: {
                display: true,
                text: 'Credit and Debit Balances',
                fontSize: 25,
            },
        },
    })
}

/**
 * @typedef Statement
 * @property {number} credit
 * @property {number} debit
 */
