import {Component, inject, OnInit} from '@angular/core';
import { UpperCasePipe} from "@angular/common";
import {LoginButtonComponent} from "../components/auth0/login-button.component";
import {MatButtonModule} from "@angular/material/button";
import {MatMenuModule} from "@angular/material/menu";
import {AuthService} from "@auth0/auth0-angular";
import {RouterLink} from "@angular/router";

@Component({
    selector: 'app-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.css'],
    imports: [UpperCasePipe, LoginButtonComponent, MatButtonModule, MatMenuModule, RouterLink]
})


export class HeaderComponent implements OnInit {

  title = "serve"
  private auth = inject(AuthService);
  constructor() { }

  ngOnInit(): void {
  }

  handleLogin(): void {
    this.auth.loginWithRedirect({
      appState: {
        target: '/profile',
      },
      authorizationParams: {
        prompt: 'login',
      },
    });
  }

}
