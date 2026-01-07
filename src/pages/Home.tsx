import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { StudentProfile, Region } from '../types';

interface HomeProps {
  profile: StudentProfile;
  setProfile: (profile: StudentProfile) => void;
}

export default function Home({ profile, setProfile }: HomeProps) {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    gpa: profile.gpa?.toString() || '',
    gpaScale: profile.gpaScale || '4.0',
    englishTest: profile.englishTest || 'none',
    englishScore: profile.englishScore?.toString() || '',
    achievements: profile.otherAchievements?.join(', ') || '',
    region: (profile.preferredRegions?.[0] as Region) || 'USA',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const newProfile: StudentProfile = {
      ...profile,
      gpa: formData.gpa ? parseFloat(formData.gpa) : undefined,
      gpaScale: formData.gpaScale as '4.0' | '5.0' | '100',
      englishTest: formData.englishTest === 'none' ? undefined : formData.englishTest as 'IELTS' | 'TOEFL',
      englishScore: formData.englishScore ? parseFloat(formData.englishScore) : undefined,
      otherAchievements: formData.achievements 
        ? formData.achievements.split(',').map(a => a.trim()).filter(a => a)
        : undefined,
      preferredRegions: [formData.region],
    };

    setProfile(newProfile);
    navigate('/results');
  };

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Узнайте свои шансы поступить в зарубежный вуз
        </h1>
        <p className="text-xl text-gray-600 max-w-2xl mx-auto">
          UniChance поможет вам быстро оценить вероятность поступления и подобрать подходящие университеты на основе вашего профиля
        </p>
      </div>

      <div className="card max-w-2xl mx-auto">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="gpa" className="block text-sm font-medium text-gray-700 mb-2">
              Средний балл (GPA)
            </label>
            <div className="flex gap-4">
              <input
                type="number"
                id="gpa"
                step="0.01"
                min="0"
                max={formData.gpaScale === '100' ? '100' : formData.gpaScale}
                value={formData.gpa}
                onChange={(e) => setFormData({ ...formData, gpa: e.target.value })}
                className="input-field flex-1"
                placeholder={formData.gpaScale === '100' ? '85.5' : '3.5'}
                required
              />
              <select
                value={formData.gpaScale}
                onChange={(e) => setFormData({ ...formData, gpaScale: e.target.value })}
                className="input-field w-32"
              >
                <option value="4.0">4.0</option>
                <option value="5.0">5.0</option>
                <option value="100">100%</option>
              </select>
            </div>
            <p className="mt-1 text-sm text-gray-500">
              Укажите ваш средний балл по выбранной шкале
            </p>
          </div>

          <div>
            <label htmlFor="englishTest" className="block text-sm font-medium text-gray-700 mb-2">
              Уровень английского
            </label>
            <div className="flex gap-4">
              <select
                id="englishTest"
                value={formData.englishTest}
                onChange={(e) => setFormData({ ...formData, englishTest: e.target.value })}
                className="input-field flex-1"
              >
                <option value="none">Не сдавал</option>
                <option value="IELTS">IELTS</option>
                <option value="TOEFL">TOEFL</option>
              </select>
              {formData.englishTest !== 'none' && (
                <input
                  type="number"
                  step="0.5"
                  min="0"
                  max={formData.englishTest === 'IELTS' ? '9' : '120'}
                  value={formData.englishScore}
                  onChange={(e) => setFormData({ ...formData, englishScore: e.target.value })}
                  className="input-field w-32"
                  placeholder={formData.englishTest === 'IELTS' ? '6.5' : '90'}
                />
              )}
            </div>
            <p className="mt-1 text-sm text-gray-500">
              {formData.englishTest === 'none' 
                ? 'Если у вас есть результаты теста, укажите их для более точной оценки'
                : `Введите ваш результат ${formData.englishTest}`
              }
            </p>
          </div>

          <div>
            <label htmlFor="achievements" className="block text-sm font-medium text-gray-700 mb-2">
              Основные достижения (опционально)
            </label>
            <input
              type="text"
              id="achievements"
              value={formData.achievements}
              onChange={(e) => setFormData({ ...formData, achievements: e.target.value })}
              className="input-field"
              placeholder="Например: призер олимпиады, волонтерство, спортивный разряд"
            />
            <p className="mt-1 text-sm text-gray-500">
              Перечислите через запятую ваши значимые достижения
            </p>
          </div>

          <div>
            <label htmlFor="region" className="block text-sm font-medium text-gray-700 mb-2">
              Желаемый регион обучения
            </label>
            <select
              id="region"
              value={formData.region}
              onChange={(e) => setFormData({ ...formData, region: e.target.value as Region })}
              className="input-field"
            >
              <option value="USA">США</option>
              <option value="UK">Великобритания</option>
              <option value="Europe">Европа</option>
              <option value="Canada">Канада</option>
              <option value="Australia">Австралия</option>
              <option value="Other">Другое</option>
            </select>
          </div>

          <button type="submit" className="btn-primary w-full text-lg py-4">
            Оценить мои шансы
          </button>
        </form>

        <div className="mt-8 pt-8 border-t border-gray-200">
          <p className="text-sm text-gray-600 text-center">
            Хотите более точную оценку?{' '}
            <a href="/profile" className="text-primary-600 hover:text-primary-700 font-medium">
              Заполните подробный профиль
            </a>
          </p>
        </div>
      </div>

      <div className="mt-12 max-w-3xl mx-auto">
        <div className="grid md:grid-cols-3 gap-6">
          <div className="text-center">
            <div className="text-3xl font-bold text-primary-600 mb-2">Быстро</div>
            <p className="text-gray-600">Оценка за несколько секунд</p>
          </div>
          <div className="text-center">
            <div className="text-3xl font-bold text-primary-600 mb-2">Точно</div>
            <p className="text-gray-600">На основе актуальных данных вузов</p>
          </div>
          <div className="text-center">
            <div className="text-3xl font-bold text-primary-600 mb-2">Полезно</div>
            <p className="text-gray-600">Персональные рекомендации по улучшению</p>
          </div>
        </div>
      </div>
    </div>
  );
}