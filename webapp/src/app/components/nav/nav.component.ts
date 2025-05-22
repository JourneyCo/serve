import { Component } from '@angular/core';
import {CommonModule} from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { Observable } from 'rxjs';
import {AuthService, HelperService} from '@services';
import { User } from '@models';
import {MaterialModule} from '@material';
import {MatToolbar} from '@angular/material/toolbar';

@Component({
  selector: 'app-nav',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MaterialModule,
    MatToolbar
  ],
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss']
})

export class NavComponent {
  isAuthenticated$: Observable<boolean>;
  user$: Observable<User | null>;
  isAdmin: Observable<boolean>;
  serve_day: Date = new Date();

  constructor(private authService: AuthService,
              private router: Router,
              private helperService: HelperService,
              ) {
    this.isAuthenticated$ = this.authService.isAuthenticated();
    this.user$ = this.authService.getCurrentUser();
    this.isAdmin = this.authService.isAdmin();
    this.serve_day = this.helperService.GetServeDate();
  }
}
