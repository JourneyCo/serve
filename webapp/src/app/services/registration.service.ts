import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class RegistrationService {
  private registrationChange = new Subject<void>();
  registrationChange$ = this.registrationChange.asObservable();

  triggerRegistrationChange() {
    this.registrationChange.next();
  }
}
