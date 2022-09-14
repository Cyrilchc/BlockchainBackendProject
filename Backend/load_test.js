import http from 'k6/http';

export let options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    stages: [
        {duration: '10s', target: 100}, // reach 100 users in 10 seconds
        {duration: '1m', target: 100}, // stays at 100 users load for 1 minutes
        {duration: '10s', target: 0}, // cool down in 10 seconds
    ],
    thresholds:{
        http_req_failed: ['rate<0.01'], // Service must be resilient
        http_req_duration: ['p(95)<300'] // Service must be fast
    }
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