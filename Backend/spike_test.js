import http from 'k6/http';

export let options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    stages: [
        {duration: '10s', target: 100}, // Normal load
        {duration: '30s', target: 1000}, // Big spike for 30 seconds
        {duration: '10s', target: 100}, // Back to normal load
        {duration: '10s', target: 0} // Recovery
    ]
};

export default () => {
    const randomUserName = Math.random().toString().substr(3, 100);
    const endpoint = 'http://localhost:5000'
    const params = {
        headers: {
            'Content-Type': 'application/json'
        }
    }

    const payload = JSON.stringify({
        username: randomUserName,
        password: 'Foobar.1',
        pin_code: '123456',
    })

    http.post(endpoint, payload, params)
}
