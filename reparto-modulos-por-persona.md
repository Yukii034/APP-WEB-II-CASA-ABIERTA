# Reparto de Módulos por Persona

## Lili — Gateway y Frontend

El **Frontend** es la interfaz visual del sistema: es donde el usuario inicia sesión, registra información, consulta recordatorios y revisa sus datos de salud. Está diseñado para que las funciones sean claras y accesibles.

El **Gateway** trabaja internamente como intermediario entre el Frontend y los demás servicios. Recibe cada solicitud, verifica a qué módulo debe enviarla y devuelve la respuesta correspondiente. Por ejemplo, cuando un usuario registra una cita, el Gateway dirige esa información al módulo de citas médicas.

## Melanie — Contacto de emergencia, cuidadores e información de salud

El módulo de **contacto de emergencia** almacena nombres, teléfonos y parentesco de las personas que deben ser notificadas en caso de una situación de riesgo o emergencia.

El módulo de **cuidadores** registra a las personas responsables del acompañamiento del usuario. Esto permite que puedan revisar información importante, como recordatorios, citas o alertas, según los permisos establecidos.

La **información de salud** reúne antecedentes médicos, alergias, enfermedades y otros datos relevantes. Tener esta información centralizada facilita una atención más rápida, organizada y segura.

## Hector — Monitoreo de signos vitales y citas médicas

El módulo de **monitoreo de signos vitales** permite registrar valores como presión arterial, frecuencia cardíaca, temperatura, oxígeno en la sangre o glucosa. El sistema guarda cada registro con su fecha y hora para consultar el historial y observar cambios.

El módulo de **citas médicas** organiza consultas, controles y exámenes. Guarda datos como fecha, hora, especialidad, médico y motivo de la cita. Además, puede utilizar esta información para generar recordatorios y evitar que el usuario olvide una atención importante.

## Macias — Recordatorio de medicamentos y alimentación

El módulo de **recordatorio de medicamentos** registra el nombre del medicamento, la dosis, la frecuencia y los horarios de consumo. Con estos datos, el sistema genera avisos para ayudar al usuario a cumplir correctamente su tratamiento.

El módulo de **alimentación** permite registrar las comidas del día, horarios, restricciones o notas nutricionales. Esto facilita el seguimiento de hábitos alimenticios y permite relacionarlos con el estado general de salud de la persona.

## Frank — Actividad física y estado de ánimo

El módulo de **actividad física** registra ejercicios como caminar, correr, realizar estiramientos o practicar deportes. Se guarda información como el tipo de actividad, duración, intensidad y fecha. Con ello, el sistema puede mostrar el progreso diario o semanal del usuario.

El módulo de **estado de ánimo** permite que el usuario indique cómo se siente, por ejemplo: feliz, tranquilo, cansado, triste o preocupado. También puede incluir una breve observación sobre lo ocurrido durante el día.

Estos módulos se relacionan porque permiten observar el bienestar físico y emocional de forma conjunta. Por ejemplo, se puede identificar si realizar actividad física frecuente se asocia con un mejor estado de ánimo.

## Winter — Reportes médicos, reportes y estimulación cognitiva

Los **reportes médicos** reúnen información de salud relevante, como signos vitales, historial de medicamentos, citas médicas y antecedentes. Esto permite contar con un resumen útil para el usuario, cuidador o profesional de salud.

El módulo general de **reportes** transforma los registros del sistema en información organizada, por ejemplo, actividades realizadas, cumplimiento de medicamentos o evolución del estado de ánimo. Así es más fácil revisar el progreso y tomar decisiones informadas.

La **estimulación cognitiva** incluye actividades orientadas a ejercitar la memoria, la atención, el razonamiento y la concentración. Estas actividades apoyan el mantenimiento de las capacidades mentales y permiten registrar el avance del usuario a lo largo del tiempo.
