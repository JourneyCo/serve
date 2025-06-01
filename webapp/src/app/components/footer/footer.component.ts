
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import {Router, RouterModule} from '@angular/router';
import { AuthService } from '@services';
import { MaterialModule } from '@material';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-footer',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MaterialModule,
  ],
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent {
  isAuthenticated$: Observable<boolean>;
  isAdmin: Observable<boolean>;

  constructor(
      private authService: AuthService,
      public router: Router
  ) {
    this.isAuthenticated$ = this.authService.isAuthenticated();
    this.isAdmin = this.authService.isAdmin();
  }

  login(): void {
    this.authService.login();
  }

  logout(): void {
    this.authService.logout();
  }
}
