import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {MapGeocoder} from '@angular/google-maps';
import {CookieService} from 'ngx-cookie-service';


@Injectable({
    providedIn: "root",
})
export class ServeCookie {
    constructor(
        private cookieService: CookieService,
    ) {}


    SetEmail(email: string) {
        this.cookieService.set("servedayemail", email, {
            expires: 180,
            secure: true,
            path: '/',
        });
    }

    DeleteEmail(email: string) {
        this.cookieService.delete("servedayemail", "/", undefined, true);
    }

    SetProject(id: number) {
        this.cookieService.set("servedayproject", String(id), {
            expires: 180,
            secure: true,
            path: '/',
        });
    }

    DeleteProject(id: number) {
        this.cookieService.delete("servedayproject", "/", undefined, true);
    }

}