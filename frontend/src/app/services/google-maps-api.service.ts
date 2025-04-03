import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { GoogleMap, MapGeocoder, MapGeocoderResponse } from '@angular/google-maps';

export interface GeocodingResult {
  latitude: number;
  longitude: number;
  formatted_address: string;
}

@Injectable({
  providedIn: 'root'
})
export class GoogleMapsApiService {
  private apiKey = environment.googleMapsApiKey;

  constructor(
    private http: HttpClient,
    private geocoder: MapGeocoder
  ) {}

  /**
   * Geocode an address using our backend service
   * @param address The address to geocode
   */
  geocodeAddress(address: string): Observable<GeocodingResult> {
    const apiUrl = `${environment.apiUrl}/api/geocode`;
    
    return this.http.post<GeocodingResult>(apiUrl, { address });
  }

  /**
   * Geocode an address using the client-side Google Maps API
   * This is the primary method that uses the Angular Google Maps library
   * @param address The address to geocode
   */
  geocodeAddressClientSide(address: string): Observable<GeocodingResult> {
    return this.geocoder.geocode({
      address: address
    }).pipe(
      map((response: MapGeocoderResponse) => {
        if (response.status !== 'OK' || !response.results.length) {
          throw new Error(`Geocoding failed: ${response.status}`);
        }
        
        const location = response.results[0].geometry.location;
        return {
          latitude: location.lat(),
          longitude: location.lng(),
          formatted_address: response.results[0].formatted_address
        };
      }),
      catchError(error => {
        console.error('Geocoding error:', error);
        return throwError(() => new Error(`Geocoding failed: ${error.message}`));
      })
    );
  }
  
  /**
   * Get Google Maps API Key
   */
  getApiKey(): string {
    return this.apiKey;
  }
}