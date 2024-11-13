import http from 'k6/http';
import { sleep } from 'k6';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '1m', target: 100 }, // Ramp up to 100 VUs in 1 minute
        { duration: '2m', target: 500 }, // Stay at 500 VUs for 2 minutes
        { duration: '1m', target: 100 }, // Ramp down to 100 VUs in 1 minute
        { duration: '1m', target: 0 },   // Ramp down to 0 VUs in 1 minute
    ],
};

export default function () {
    let res = http.get('http://localhost:4000');
    check(res, {
        'status was 200': (r) => r.status === 200,
        'transaction time OK': (r) => r.timings.duration < 200,
    });
    sleep(1); // Sleep for 1 second between requests
}
