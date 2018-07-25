import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectsCclaSelectModal } from './projects-ccla-select-modal';

@NgModule({
  declarations: [
    ProjectsCclaSelectModal,
  ],
  imports: [
    IonicPageModule.forChild(ProjectsCclaSelectModal)
  ],
  entryComponents: [
    ProjectsCclaSelectModal,
  ]
})
export class ProjectsCclaSelectModalModule {}