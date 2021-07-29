create table `audio` (
 `id` integer primary key autoincrement
,`lang` varchar(5)
,`text` varchar(100) not null
,`audio_content` text not null
,unique(`lang`, `text`)
);
