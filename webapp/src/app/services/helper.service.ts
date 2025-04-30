import {Injectable} from '@angular/core';
import {MatSnackBar} from '@angular/material/snack-bar';

@Injectable({
  providedIn: "root",
})
export class HelperService {

  constructor(
    private snackBar: MatSnackBar,
  ) {
  }


  GetServeDate(): Date {
    const serve_day = import.meta.env['NG_APP_SERVE_DAY'];
    const [yearStr, monthStr, dayStr] = serve_day.split("-");
    const year = parseInt(yearStr, 10); // Assumes 21st century for 2-digit years
    const month = parseInt(monthStr, 10) - 1; // JavaScript months are 0-indexed
    const day = parseInt(dayStr, 10);

    return new Date(year, month, day);
  }

  showSuccess(message: string): void {
    this.snackBar.open(message, "Close", {
      duration: 3000,
      panelClass: ["success-snackbar"],
    });
  }

  showError(message: string): void {
    this.snackBar.open(message, "Close", {
      duration: 5000,
      panelClass: ["error-snackbar"],
    });
  }

}
