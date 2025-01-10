import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {RouterModule } from '@angular/router';
import { ROUTES } from "./routes"

@NgModule({
  declarations: [],
  imports: [
    CommonModule, RouterModule.forRoot(ROUTES)
  ],
  exports: [RouterModule]
})
export class AppRoutingModule { }
