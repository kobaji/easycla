# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

from unittest.mock import patch, Mock

import pytest

from cla.models.dynamo_models import User, Project, Company, CCLAWhitelistRequest
from cla.models.event_types import EventType
from cla.controllers import user as user_controller
from cla.auth import AuthUser


@pytest.fixture
def create_event_user():
    user_controller.create_event = Mock()


@patch('cla.controllers.user.Event.create_event')
def test_request_company_whitelist(mock_event, create_event_user, project, company, user):
    """ Test user requesting to be added to the Approved List event """
    event_type = EventType.RequestCompanyWL
    User.load = Mock()
    User.get_user_name = Mock(return_value=user.get_user_name())
    User.get_user_emails = Mock(return_value=[user.get_user_email()])
    User.get_user_email = Mock(return_value=user.get_user_email())
    Company.load = Mock()
    Company.get_company_name = Mock(return_value=company.get_company_name())
    Project.load = Mock()
    Project.get_project_name = Mock(return_value=project.get_project_name())
    user_controller.get_email_service = Mock()
    user_controller.send = Mock()
    user_controller.request_company_whitelist(
        user.get_user_id(),
        company.get_company_id(),
        user.get_user_name(),
        user.get_user_email(),
        project.get_project_id(),
        message="Please add",
        recipient_name="Recipient Name",
        recipient_email="Recipient Email",
    )

    event_data = (f'CLA: contributor {user.get_user_name()} requests to be Approved for the '
                  f'project: {project.get_project_name()} '
                  f'organization: {company.get_company_name()} '
                  f'as {user.get_user_name()} <{user.get_user_email()}>')

    mock_event.assert_called_once_with(
       event_user_id=user.get_user_id(),
       event_project_id=project.get_project_id(),
       event_company_id=company.get_company_id(),
       event_type=event_type,
       event_data=event_data,
       event_summary=event_data,
       contains_pii=True,
    )


@patch('cla.controllers.user.Event.create_event')
def test_invite_cla_manager(mock_event, create_event_user, user):
    """ Test send email to CLA manager event """
    User.load = Mock()
    CCLAWhitelistRequest.save = Mock()
    user_controller.send_email_to_cla_manager = Mock()
    contributor_id = user.get_user_id()
    contributor_name = user.get_user_name()
    contributor_email = user.get_user_email()
    cla_manager_name = "admin"
    cla_manager_email = "foo@admin.com"
    project_name = "foo_project"
    company_name = "Test Company"
    event_data = (f'sent email to CLA Manager: {cla_manager_name} with email {cla_manager_email} '
                  f'for project {project_name} and company {company_name} '
                  f'to user {contributor_name} with email {contributor_email}')
    # TODO FIX Unit test - need to mock Project load_project_by_name() function
    #user_controller.invite_cla_manager(contributor_id, contributor_name, contributor_email,
    #                                   cla_manager_name, cla_manager_email,
    #                                   project_name, company_name)
    #mock_event.assert_called_once_with(
    #    event_user_id=contributor_id,
    #    event_project_name=project_name,
    #    event_data=event_data,
    #    event_type=EventType.InviteAdmin,
    #    contains_pii=True,
    #)


@patch('cla.controllers.user.Event.create_event')
def test_request_company_ccla(mock_event, create_event_user, user, project, company):
    """ Test request company ccla event """
    User.load = Mock()
    User.get_user_name = Mock(return_value=user.get_user_name())
    email = user.get_user_email()
    Company.load = Mock()
    Project.load = Mock()
    Project.get_project_name = Mock(return_value=project.get_project_name())
    manager = User(lf_username="harold", user_email="foo@gmail.com")
    Company.get_managers = Mock(return_value=[manager, ])
    event_data = f"Sent email to sign ccla for {project.get_project_name() }"
    user_controller.request_company_ccla(
        user.get_user_id(), email, company.get_company_id(), project.get_project_id()
    )
    mock_event.assert_called_once_with(
        event_data=event_data,
        event_summary=event_data,
        event_type=EventType.RequestCCLA,
        event_user_id=user.get_user_id(),
        event_company_id=company.get_company_id(),
        contains_pii=False,
    )
