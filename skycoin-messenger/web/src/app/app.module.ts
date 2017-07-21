import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { AppComponent } from './app.component';
import {
  ImViewComponent,
  ImRecentBarComponent,
  ImRecentItemComponent,
  ImHeadComponent,
  ImHistoryViewComponent,
  ImHistoryMessageComponent,
  CreateChatDialogComponent,
  AlertDialogComponent
} from '../components';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SocketService, ToolService } from '../providers';
import {
  MdCheckboxModule,
  MdMenuModule,
  MdIconModule,
  MdDialogModule,
  MdInputModule,
  MdButtonModule,
} from '@angular/material';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import 'web-animations-js'


@NgModule({
  declarations: [
    AppComponent,
    ImViewComponent,
    ImRecentBarComponent,
    ImRecentItemComponent,
    ImHeadComponent,
    ImHistoryViewComponent,
    ImHistoryMessageComponent,
    CreateChatDialogComponent,
    AlertDialogComponent
  ],
  imports: [
    FormsModule,
    ReactiveFormsModule,
    BrowserModule,
    FlexLayoutModule,
    BrowserAnimationsModule,
    MdCheckboxModule,
    MdMenuModule,
    MdIconModule,
    MdDialogModule,
    MdInputModule,
    MdButtonModule,
  ],
  entryComponents: [
    CreateChatDialogComponent,
    AlertDialogComponent
  ],
  providers: [SocketService, ToolService],
  bootstrap: [AppComponent]
})
export class AppModule { }
