import {Component, OnInit, ViewChild} from '@angular/core';
import {GoogleMap, MapAdvancedMarker, MapAnchorPoint, MapInfoWindow} from "@angular/google-maps";
import {CommonModule} from "@angular/common";
import {APIService} from "@services";
import {Location, Project} from "@models";
import {forkJoin} from "rxjs";

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrl: './map.component.css',
  standalone: true,
  imports: [CommonModule, GoogleMap, MapInfoWindow, MapAdvancedMarker]
})
export class MapComponent implements OnInit {
  @ViewChild(MapInfoWindow, { static: false })
  infoWindow!: MapInfoWindow
  markers: any[] = [];
  selectedMarker: any;

  options: google.maps.MapOptions = {
    center: {lat: 39.491482, lng: -104.874878},
    zoom: 12,
    mapId: '18c474b41c1ac65a',
  };
  display: google.maps.LatLngLiteral = {lat: 39.491482, lng: -104.874878};


  constructor(private APIService: APIService) {}

  ngOnInit() {

    const observable = forkJoin({
      projects: this.APIService.getProjects(),
      locations: this.APIService.getLocations(),
    });

    observable.subscribe(data => {
      data.locations.forEach((location: Location)=> {
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
        data.projects.forEach((project: Project) => {
          if (project.location_id == location.id) {
            m.project = project;
          }
        });
        this.markers.push(m);
      });
    });
  }

  moveMap(event: google.maps.MapMouseEvent) {
    // @ts-ignore
    this.center = event.latLng.toJSON();
  }

  move(event: google.maps.MapMouseEvent) {
    // @ts-ignore
    this.display = event.latLng.toJSON();
  }

  openInfoWindow(marker: MapAdvancedMarker, m: any) {
    this.selectedMarker = m;
    this.infoWindow.open(marker);
  }
}
