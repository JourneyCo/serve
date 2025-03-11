import {Component, inject, OnInit} from '@angular/core';
import {MatButtonModule} from "@angular/material/button";
import {MatMenuModule} from "@angular/material/menu";
import {AuthService} from "@auth0/auth0-angular";
import {RouterLink} from "@angular/router";
import {DOCUMENT} from "@angular/common";

@Component({
    selector: 'app-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.css'],
    imports: [MatButtonModule, MatMenuModule, RouterLink]
})


export class HeaderComponent implements OnInit {

  title = "serve"
  private auth = inject(AuthService);
  private doc = inject(DOCUMENT)
  constructor() { }

  ngOnInit(): void {
  }

  handleLogin(): void {
    this.auth.loginWithRedirect({
      appState: {
        target: '/projects',
      },
      authorizationParams: {
        prompt: 'login',
      },
    });
  }

  handleLogout(): void {
    this.auth.logout({
      logoutParams: {
        returnTo: this.doc.location.origin,
      },
    });
  }

}
