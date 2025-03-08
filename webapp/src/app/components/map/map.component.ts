import {Component, Input, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {GoogleMap, MapAdvancedMarker, MapInfoWindow} from "@angular/google-maps";
import {CommonModule} from "@angular/common";
import {Location, Project} from "@models";
import {Observable, Subscription} from "rxjs";

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrl: './map.component.css',
  standalone: true,
  imports: [CommonModule, GoogleMap, MapInfoWindow, MapAdvancedMarker]
})
export class MapComponent implements OnInit, OnDestroy {
  @ViewChild(MapInfoWindow, { static: false }) infoWindow!: MapInfoWindow;
  @ViewChild('mapRef') mapRef: GoogleMap;
  @Input() locations: any[];
  @Input() projects: any[];
  @Input() events: Observable<void>;
  private eventsSubscription: Subscription

  center: google.maps.LatLngLiteral = {lat: 39.491482, lng: -104.874878};
  markers: any[] = [];
  selectedMarker: any;
  zoom = 12;

  options: google.maps.MapOptions = {
    mapId: '18c474b41c1ac65a',
  };
  display: google.maps.LatLngLiteral = {lat: 39.491482, lng: -104.874878};


  constructor() {}

  ngOnInit() {
    this.eventsSubscription = this.events?.subscribe((data) => {
      this.moveToClickedRow(data);
      }
    );
    this.locations?.forEach((location: Location)=> {
      let loc: google.maps.LatLngLiteral = {
        lat: location.latitude,
        lng: location.longitude
      }
      // @ts-ignore
      let mark: google.maps.marker.AdvancedMarkerElement = {
        position: loc,
        title: location.info,
      }
      let m: any = {
        marker: mark,
        info: location.info
      }
      this.projects.forEach((project: Project) => {
        if (project.location_id == location.id) {
          m.project = project;
        }
      });
      this.markers.push(m);
    });
  }

  ngOnDestroy() {
    if (this.eventsSubscription) {
      this.eventsSubscription.unsubscribe();
    }
  }

  moveMap(event: google.maps.MapMouseEvent) {
    this.center = event.latLng!.toJSON();
  }

  moveToClickedRow(location: any ) {
    const l = new google.maps.LatLng(location.latitude, location.longitude);
    this.center = l.toJSON();
  }

  move(event: google.maps.MapMouseEvent) {
    this.display = event.latLng!.toJSON();
  }

  openInfoWindow(marker: MapAdvancedMarker, m: any) {
    this.selectedMarker = m;
    this.infoWindow.open(marker);
  }
}
