import { Component } from "@angular/core";
import { RouterOutlet } from "@angular/router";
import { HttpClientModule } from "@angular/common/http";
import { CommonModule } from "@angular/common";
import { NavComponent } from '@components';

@Component({
  selector: "app-root",
  standalone: true,
  imports: [CommonModule, RouterOutlet, HttpClientModule, NavComponent],
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"],
})
export class AppComponent {
  title = "Journey Serve Day";

  constructor() {}
}
