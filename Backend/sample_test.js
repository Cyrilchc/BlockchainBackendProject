import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    vus: 1,
    duration: '10s'
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
        username : randomUserName,
        password : 'Foobar.1',
        pin_code : '123456',
    })

    http.post(endpoint, payload, params)
    sleep(1);
}