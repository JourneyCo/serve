import {Component, inject, OnDestroy, OnInit} from '@angular/core';
import {MatButtonModule} from "@angular/material/button";
import {MatMenuModule} from "@angular/material/menu";
import {AuthService} from "@auth0/auth0-angular";
import {RouterLink} from "@angular/router";
import {CommonModule, DOCUMENT} from "@angular/common";
import {Subscription} from "rxjs";

@Component({
    selector: 'app-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.css'],
    imports: [CommonModule, MatButtonModule, MatMenuModule, RouterLink]
})


export class HeaderComponent implements OnInit, OnDestroy {

  title = "serve"
  private auth = inject(AuthService);
  private doc = inject(DOCUMENT)
  isLoggedIn: boolean = false;
  private authSubscription?: Subscription;

  constructor(public authService: AuthService) { }

  ngOnInit(): void {
    this.authSubscription = this.auth.isAuthenticated$.subscribe(
      (isAuthenticated) => {
        this.isLoggedIn = isAuthenticated;
      }
    );
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

  ngOnDestroy(): void {
    if (this.authSubscription) {
      this.authSubscription.unsubscribe();
    }
  }

}
