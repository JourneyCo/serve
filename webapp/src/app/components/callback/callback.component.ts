import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
    selector: 'app-callback',
    template: '<p>Redirecting...</p>'
})
export class CallbackComponent implements OnInit {
    constructor(private router: Router) {}

    ngOnInit(): void {
        const redirectUrl = localStorage.getItem('redirectUrl');
        if (redirectUrl) {
            localStorage.removeItem('redirectUrl'); // Clean up
            this.router.navigateByUrl(redirectUrl);
        } else {
            this.router.navigate(['/']); // Fallback if no URL is stored
        }
    }
}