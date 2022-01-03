create table `recordbook` (
 `app_user_id` int not null
,`workbook_id` int not null
,`problem_type_id` int not null
,`study_type_id` int not null
,`problem_id` int not null
,`result_prev3` tinyint
,`result_prev2` tinyint
,`result_prev1` tinyint
,`level` int not null
,`memorized` tinyint
,`last_answered_at` datetime not null default current_timestamp on update current_timestamp
,primary key(`app_user_id`, `problem_id`, `study_type_id`, `problem_type_id`)
,foreign key(`app_user_id`) references `app_user`(`id`) on delete cascade
,foreign key(`problem_type_id`) references `problem_type`(`id`) on delete cascade
,foreign key(`study_type_id`) references `study_type`(`id`) on delete cascade
,foreign key(`workbook_id`) references `workbook`(`id`) on delete cascade
);
