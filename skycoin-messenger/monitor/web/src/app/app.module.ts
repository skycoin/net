import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {
  MdGridListModule,
  MatListModule,
  MdIconModule,
  MatTableModule,
  MatTooltipModule,
  MatChipsModule,
  MatSnackBarModule
} from '@angular/material';
import { AppComponent } from './app.component';
import { HttpClientModule } from '@angular/common/http';
import { ApiService } from './service';
import { TimeAgoPipe, ByteToPipe } from './pipe';
import { DashboardComponent } from './page';

import { AppRoutingModule } from './route/app-routing.module';

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    TimeAgoPipe,
    ByteToPipe
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    MdGridListModule,
    MatListModule,
    MdIconModule,
    MatTableModule,
    MatTooltipModule,
    MatChipsModule,
    MatSnackBarModule
  ],
  providers: [ApiService],
  bootstrap: [AppComponent]
})
export class AppModule { }
