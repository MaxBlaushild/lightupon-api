INSERT INTO trips
(title, description, image_url)
VALUES
('First Trip (HQ)', 'The Very First One Hombre', 'http://i.imgur.com/eMXDLSP.jpg'),
('Second Trip (HQ)', 'The Very Second One Hombre', 'http://i.imgur.com/eMXDLSP.jpg'),
('Copley Trip', 'This one leaves from Copley', 'http://i.imgur.com/eMXDLSP.jpg');


INSERT INTO scenes
(name, latitude, longitude, trip_id, scene_order)
VALUES
('Headquarters', 42.338651, -71.079991, 1, 1),
('Laced', 42.340940, -71.082150, 1, 2),
('The Church Park', 42.343441, -71.079840, 1, 3),
('place 4', 42.358651, -71.073991, 1, 4),
('place 5', 42.360940, -71.081150, 1, 5),
('place 6', 42.373441, -71.078840, 1, 6);


INSERT INTO cards
(scene_id, card_order, dialogue, universal)
VALUES
(1, 1, 'Good Morning! Thank God youre awake!', true),
(2, 1, 'The sale is almost half over!', true),
(3, 1, 'If you hurry, youl still finds some good deals!', true),
(4, 1, 'Plus the friendliest staff around. Always willing. Always ready.', true),
(5, 1, 'If you buy 2 or more liters, you get a 3rd free! Buy a gallon, and receive half off your next gallon and a free breakfast sandwich, made fresh daily', true),
(6, 1, 'Cumberland Farms, putting the Cum in Berland Farms', true);


INSERT INTO users
(facebook_id, email, token)
VALUES
('1145256482154055', 'campbelldaley@gmail.com', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTc5NzEwNTUsImZhY2Vib29rSWQiOiIxMTQ1MjU2NDgyMTU0MDU1In0.ji2VyJmDuxiBnBYd19gGqvb7GzAGoVBf0lngZ3UzceA'),
('1345256481024053', 'max@max.com', 'osdfgads8sduher9eso9.eyJleHAiOjE0NTc1NjcwNjAsImZhY2Vib29rSWQiOiIxMTQ1MjU2NDgyMTU0MDU1In0.2Apc5F9yJL3DmlivvrlalvTRmHG7jZ8sh0l_iUqrYyU');


INSERT INTO partyusers
(party_id, user_id)
VALUES
(1, 2);



