import { Component, OnInit } from '@angular/core';
import {CommonModule, DatePipe} from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { Observable } from 'rxjs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatDividerModule } from '@angular/material/divider';
import { AuthService } from '@services';
import { User } from '@models';
import {environment} from "../../../environments/environment";

@Component({
  selector: 'app-nav',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule,
    MatDividerModule,
      DatePipe,
  ],
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss']
})
export class NavComponent implements OnInit {
  isAuthenticated$: Observable<boolean>;
  user$: Observable<User | null>;
  isAdmin: Observable<boolean>;
  serve_day: Date = new Date();

  constructor(private authService: AuthService,
              private router: Router,
              ) {
    this.isAuthenticated$ = this.authService.isAuthenticated();
    this.user$ = this.authService.getCurrentUser();
    this.isAdmin = this.authService.isAdmin();
    this.getServeDay();
  }

  ngOnInit(): void {}

  login(): void {
    this.authService.login();
  }

  logout(): void {
    this.authService.logout();
  }

  getServeDay() : void {
    const serve = environment.serveDay;
    const split = serve.split("-");
    this.serve_day.setMonth(Number(split[0]) - 1);
    this.serve_day.setDate(Number(split[1]));
  }
}