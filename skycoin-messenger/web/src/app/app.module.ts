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
  AlertDialogComponent,
  ImInfoDialogComponent
} from '../components';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SocketService, UserService } from '../providers';
import { ToolService } from '../providers/tool/tool.service';
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
    AlertDialogComponent,
    ImInfoDialogComponent
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
    AlertDialogComponent,
    ImInfoDialogComponent
  ],
  providers: [SocketService, UserService, ToolService],
  bootstrap: [AppComponent]
})
export class AppModule { }
