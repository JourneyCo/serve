import {Injectable} from '@angular/core';

@Injectable({
  providedIn: "root",
})
export class HelperService {


  GetServeDate(): Date {
    const serve_day = import.meta.env['NG_APP_SERVE_DAY'];
    const [monthStr, dayStr, yearStr] = serve_day.split("-");
    const year = 2000 + parseInt(yearStr, 10); // Assumes 21st century for 2-digit years
    const month = parseInt(monthStr, 10) - 1; // JavaScript months are 0-indexed
    const day = parseInt(dayStr, 10);

    return new Date(year, month, day);
  }

}
